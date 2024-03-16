package ngrams

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	// Check if this is a zip file
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".zip" {
		return ft.processZipFile(ctx, path, fn)
	}

	// Normal file
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

func (ft *FrequencyTable) processZipFile(ctx context.Context, path string, fn processFileFunc) error {
	zf, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %q. %w", path, err)
	}
	defer func() {
		if err := zf.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close zip file %s. %v", path, err)
		}
	}()

	closer := func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close file. %v", err)
		}
	}

	filter := func(f *zip.File) bool {
		// Ignore directories
		if f.FileInfo().IsDir() {
			return false
		}

		// Ignore hidden files (especially pesky .DS_Store)
		if strings.HasPrefix(filepath.Base(f.Name), ".") {
			return false
		}

		return true
	}

	// Get more up to date progress size
	totalUncompressedSize := uint64(0)
	for _, f := range zf.File {
		if !filter(f) {
			continue
		}
		totalUncompressedSize += f.UncompressedSize64
	}
	ft.progressReporter.AddToTotalSize(totalUncompressedSize)

	// Process each file in the zip
	for _, f := range zf.File {
		if !filter(f) {
			continue
		}

		//AJ### TODO: Need to look at progress reporting, need to update based on uncompressed size etc
		zr, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %q inside of zip file %q. %w", f.Name, path, err)
		}

		r := bufio.NewReader(zr)
		err = fn(ctx, r)
		if err != nil {
			closer(zr)
			return err
		}
		closer(zr)
	}

	return nil
}
