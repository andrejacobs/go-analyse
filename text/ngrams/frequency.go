package ngrams

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/andrejacobs/go-collection/collection"
	"golang.org/x/exp/maps"
)

type FrequencyTable struct {
	frequencies tokenFrequencyMap
}

// LoadFrequencies parses a frequency table from an io.Reader.
//
// Expected CSV format in UTF-8: token,count,percentage
// Lines starting with a # is ignored.
func LoadFrequencies(r io.Reader) (FrequencyTable, error) {
	result := FrequencyTable{
		frequencies: make(tokenFrequencyMap),
	}
	csvR := csv.NewReader(r)

	for {
		record, err := csvR.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return FrequencyTable{}, fmt.Errorf("failed to parse csv. %w", err)
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
			return FrequencyTable{}, fmt.Errorf("failed to parse the count field from the csv. %v. %w", record, err)
		}

		percentage, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 32)
		if err != nil {
			return FrequencyTable{}, fmt.Errorf("failed to parse the percentage field from the csv. %v. %w", record, err)
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
func LoadFrequenciesFromFile(path string) (FrequencyTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return FrequencyTable{}, fmt.Errorf("failed to open %s. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close %s. %v", path, err)
		}
	}()

	result, err := LoadFrequencies(f)
	if err != nil {
		return FrequencyTable{}, fmt.Errorf("failed to load languages from %q. %w", path, err)
	}

	return result, nil
}

// Len returns the number of [Frequency] entries in the table.
func (ft FrequencyTable) Len() int {
	return len(ft.frequencies)
}

// Entries returns the token frequencies in the table.
// NOTE: The order can not be guarenteed since the underlying data structure uses a map.
func (ft FrequencyTable) Entries() []Frequency {
	return maps.Values(ft.frequencies)
}

// EntriesSortedByCount returns the token frequencies in the table sorted by the count (descending) going from
// the token that appears the most to the least (highest to lowest frequency).
func (ft FrequencyTable) EntriesSortedByCount() []Frequency {
	kvs := collection.MapSortedByValueFunc(ft.frequencies, func(l, r Frequency) bool {
		return l.Count > r.Count
	})

	values := collection.JustValues(kvs)
	slices.SortFunc(values, func(a Frequency, b Frequency) int {
		if a.Count > b.Count {
			return -1
		} else if a.Count < b.Count {
			return 1
		} else {
			return 0
		}
	})

	return values
}

// Tokens returns the unique tokens present in the table.
// NOTE: The order can not be guarenteed since the underlying data structure uses a map.
func (ft FrequencyTable) Tokens() []string {
	return maps.Keys(ft.frequencies)
}

// Get returns the frequency information for the given token.
// A bool is also returned to indicate if the token does exist in the table or not.
func (ft FrequencyTable) Get(token string) (Frequency, bool) {
	result, exists := ft.frequencies[token]
	return result, exists
}

//TODO: add a save function
// add token
// calculate freqs

//-----------------------------------------------------------------------------

type Frequency struct {
	Token      string
	Count      int
	Percentage float32
}

type tokenFrequencyMap map[string]Frequency
