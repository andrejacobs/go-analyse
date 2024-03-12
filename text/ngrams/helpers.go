package ngrams

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/andrejacobs/go-analyse/text/alphabet"
)

// FrequencyTableByParsingLetters creates a new frequency table by parsing letter ngrams from the given input paths.
func FrequencyTableByParsingLetters(ctx context.Context, paths []string,
	language alphabet.Language, tokenSize int) (*FrequencyTable, error) {

	ft := NewFrequencyTable()

	for _, path := range paths {
		if err := ft.parseLettersFromFile(ctx, path, language, tokenSize); err != nil {
			return nil, fmt.Errorf("failed to create the frequency table from the file %q. %w", path, err)
		}
	}

	ft.Update()
	return ft, nil
}

func (ft *FrequencyTable) parseLettersFromFile(ctx context.Context, path string,
	language alphabet.Language, tokenSize int) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to parse the frequency table from the file %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close %s. %v", path, err)
		}
	}()

	r := bufio.NewReader(f)

	err = ft.ParseLetterTokens(ctx, r, language, tokenSize)
	if err != nil {
		return fmt.Errorf("failed to parse the frequency table from the file %q. %w", path, err)
	}

	return nil
}
