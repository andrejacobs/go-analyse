package ngrams

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/andrejacobs/go-analyse/text/alphabet"
)

// FrequencyTableByParsingLetters creates a new frequency table by parsing letter ngrams from the given text file.
func FrequencyTableByParsingLetters(ctx context.Context, path string,
	language alphabet.Language, tokenSize int) (*FrequencyTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create the frequency table from the file %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close %s. %v", path, err)
		}
	}()

	ft := NewFrequencyTable()
	r := bufio.NewReader(f)

	err = ft.ParseLetterTokens(ctx, r, language, tokenSize)
	if err != nil {
		return ft, fmt.Errorf("failed to create the frequency table from the file %q. %w", path, err)
	}
	return ft, nil
}
