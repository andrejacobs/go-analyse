package ngrams

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/andrejacobs/go-analyse/text/alphabet"
)

// UpdateTableByParsingLettersFromFiles updates the frequency table by parsing letter ngrams from the given input paths.
func (ft *FrequencyTable) UpdateTableByParsingLettersFromFiles(ctx context.Context, paths []string,
	language alphabet.Language, tokenSize int) error {

	total := len(paths)

	for i, path := range paths {
		ft.progressReporter.Started(path, i, total)
		if err := ft.parseLettersFromFile(ctx, path, language, tokenSize); err != nil {
			return fmt.Errorf("failed to update the frequency table from the file %q. %w", path, err)
		}
	}

	ft.Update()
	return nil
}

// UpdateTableByParsingWordsFromFiles updates the frequency table by parsing word ngrams from the given input paths.
func (ft *FrequencyTable) UpdateTableByParsingWordsFromFiles(ctx context.Context, paths []string,
	language alphabet.Language, tokenSize int) error {

	total := len(paths)

	for i, path := range paths {
		ft.progressReporter.Started(path, i, total)
		if err := ft.parseWordsFromFile(ctx, path, language, tokenSize); err != nil {
			return fmt.Errorf("failed to update the frequency table from the file %q. %w", path, err)
		}
	}

	ft.Update()
	return nil
}

func (ft *FrequencyTable) parseLettersFromFile(ctx context.Context, path string,
	language alphabet.Language, tokenSize int) error {
	return ft.processFile(ctx, path, func(ctx context.Context, r io.Reader) error {
		err := ft.ParseLetterTokens(ctx, r, language, tokenSize)
		if err != nil {
			return fmt.Errorf("failed to parse the frequency table from the file %q. %w", path, err)
		}
		return nil
	})
}

func (ft *FrequencyTable) parseWordsFromFile(ctx context.Context, path string,
	language alphabet.Language, tokenSize int) error {
	return ft.processFile(ctx, path, func(ctx context.Context, r io.Reader) error {
		err := ft.ParseWordTokens(ctx, r, language, tokenSize)
		if err != nil {
			return fmt.Errorf("failed to parse the frequency table from the file %q. %w", path, err)
		}
		return nil
	})
}

type processFileFunc func(ctx context.Context, r io.Reader) error

func (ft *FrequencyTable) processFile(ctx context.Context, path string, fn processFileFunc) error {
	//AJ### TODO: Add zip support

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open the file %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close %s. %v", path, err)
		}
	}()

	r := bufio.NewReader(f)
	err = fn(ctx, r)
	if err != nil {
		return err
	}
	return nil
}
