package ngrams

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-collection/collection"
	"golang.org/x/exp/maps"
)

type FrequencyTable struct {
	frequencies tokenFrequencyMap
	mu          sync.RWMutex
}

// LoadFrequencies parses a frequency table from an io.Reader.
//
// Expected CSV format in UTF-8: token,count,percentage
// Lines starting with a # is ignored.
func LoadFrequencies(r io.Reader) (*FrequencyTable, error) {
	result := &FrequencyTable{
		frequencies: make(tokenFrequencyMap),
	}
	csvR := csv.NewReader(r)

	for {
		record, err := csvR.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse csv. %w", err)
		}

		if len(record) < 3 {
			continue
		}

		if strings.HasPrefix(record[0], "#") {
			continue
		}

		token := record[0]

		count, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			return nil, fmt.Errorf("failed to parse the count field from the csv. %v. %w", record, err)
		}

		percentage, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the percentage field from the csv. %v. %w", record, err)
		}

		freq := Frequency{
			Token:      token,
			Count:      count,
			Percentage: float32(percentage),
		}

		result.frequencies[token] = freq
	}

	return result, nil
}

// Load a set of languages from a UTF-8 encoded text file.
// See LoadLanguages for more details.
func LoadFrequenciesFromFile(path string) (*FrequencyTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close %s. %v", path, err)
		}
	}()

	result, err := LoadFrequencies(f)
	if err != nil {
		return nil, fmt.Errorf("failed to load languages from %q. %w", path, err)
	}

	return result, nil
}

// NewFrequencyTable creates a new [FrequencyTable].
func NewFrequencyTable() *FrequencyTable {
	return &FrequencyTable{
		frequencies: make(tokenFrequencyMap),
	}
}

// Len returns the number of [Frequency] entries in the table.
func (ft *FrequencyTable) Len() int {
	ft.mu.RLock()
	defer ft.mu.RUnlock()
	return len(ft.frequencies)
}

// Entries returns the token frequencies in the table.
// NOTE: The order can not be guarenteed since the underlying data structure uses a map.
func (ft *FrequencyTable) Entries() []Frequency {
	ft.mu.RLock()
	defer ft.mu.RUnlock()
	return maps.Values(ft.frequencies)
}

// EntriesSortedByCount returns the token frequencies in the table sorted by the count (descending) going from
// the token that appears the most to the least (highest to lowest frequency).
func (ft *FrequencyTable) EntriesSortedByCount() []Frequency {
	ft.mu.RLock()
	kvs := collection.MapSortedByValueFunc(ft.frequencies, func(l, r Frequency) bool {
		return l.Count > r.Count
	})
	ft.mu.RUnlock()

	values := collection.JustValues(kvs)
	slices.SortFunc(values, func(a Frequency, b Frequency) int {
		if a.Count > b.Count {
			return -1
		} else if a.Count < b.Count {
			return 1
		} else {
			return strings.Compare(a.Token, b.Token)
		}
	})

	return values
}

// Tokens returns the unique tokens present in the table.
// NOTE: The order can not be guarenteed since the underlying data structure uses a map.
func (ft *FrequencyTable) Tokens() []string {
	ft.mu.RLock()
	defer ft.mu.RUnlock()
	return maps.Keys(ft.frequencies)
}

// Get returns the frequency information for the given token.
// A bool is also returned to indicate if the token does exist in the table or not.
func (ft *FrequencyTable) Get(token string) (Frequency, bool) {
	ft.mu.RLock()
	defer ft.mu.RUnlock()
	result, exists := ft.frequencies[token]
	return result, exists
}

// Add a token with the given frequency count.
// If the token has already been added then it's count will be incremented.
func (ft *FrequencyTable) Add(token string, count int) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	freq, exists := ft.frequencies[token]
	if !exists {
		ft.frequencies[token] = Frequency{Token: token, Count: count}
	} else {
		freq.Count += count
		ft.frequencies[token] = freq
	}
}

// Save the frequency table to the io.Writer in the same CSV format used by the Load functions.
func (ft *FrequencyTable) Save(w io.Writer) error {
	csvW := csv.NewWriter(w)
	err := csvW.Write([]string{"#token", "count", "percentage"})
	if err != nil {
		return fmt.Errorf("failed to write the csv header. %w", err)
	}

	freqs := ft.EntriesSortedByCount()
	for _, freq := range freqs {
		err := csvW.Write([]string{freq.Token, strconv.Itoa(freq.Count), strconv.FormatFloat(float64(freq.Percentage), 'f', 8, 32)})
		if err != nil {
			return fmt.Errorf("failed to write the token %q. %w", freq.Token, err)
		}
	}

	csvW.Flush()
	if err := csvW.Error(); err != nil {
		return fmt.Errorf("failed to write the frequency table. %w", err)
	}
	return nil
}

// Update will calculate and update the token frequencies
func (ft *FrequencyTable) Update() {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	sum := 0
	for _, freq := range ft.frequencies {
		sum += freq.Count
	}

	for k, freq := range ft.frequencies {
		freq.Percentage = float32(freq.Count) / float32(sum)
		ft.frequencies[k] = freq
	}
}

//-----------------------------------------------------------------------------

// ParseLetterTokens is used to parse ngrams for letter combinations of the given tokenSize and language
// from the io.Reader and then update the frequency table.
func (ft *FrequencyTable) ParseLetterTokens(ctx context.Context, input io.Reader, language alphabet.Language,
	tokenSize int) error {

	err := ParseLetterTokens(ctx, input, language, tokenSize,
		func(token string, err error) error {
			if err == nil {
				ft.Add(token, 1)
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("failed to parse the letter tokens. %w", err)
	}
	return nil
}

// ParseWordTokens is used to parse ngrams for word combinations of the given tokenSize and language
// from the io.Reader and then update the frequency table.
func (ft *FrequencyTable) ParseWordTokens(ctx context.Context, input io.Reader, language alphabet.Language,
	tokenSize int) error {

	err := ParseWordTokens(ctx, input, language, tokenSize,
		func(token string, err error) error {
			if err == nil {
				ft.Add(token, 1)
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("failed to parse the word tokens. %w", err)
	}
	return nil
}

//-----------------------------------------------------------------------------

type Frequency struct {
	Token      string
	Count      int
	Percentage float32
}

type tokenFrequencyMap map[string]Frequency
