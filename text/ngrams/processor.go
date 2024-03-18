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

// Processor is used to parse letter or word ngrams from input sources.
type Processor struct {
	ft               *FrequencyTable
	progressReporter ProgressReporter
	language         alphabet.Language
	tokenSize        int
	mode             ProcessorMode
}

// ProcessorMode specifies whether the processor works on letter or word ngrams
type ProcessorMode bool

const (
	ProcessLetters ProcessorMode = false
	ProcessWords   ProcessorMode = true
)

// NewProcessor creates a new frequency table and does not report progress.
func NewProcessor(mode ProcessorMode, language alphabet.Language, tokenSize int) *Processor {
	p := &Processor{
		ft:               NewFrequencyTable(),
		progressReporter: &nullProgressReporter{},
		language:         language,
		tokenSize:        tokenSize,
		mode:             mode,
	}
	return p
}

// SetProgressReporter sets the progress reporter to use.
func (p *Processor) SetProgressReporter(reporter ProgressReporter) {
	p.progressReporter = reporter
}

// Table returns the frequency table.
func (p *Processor) FrequencyTable() *FrequencyTable {
	return p.ft
}

// LoadFrequenciesFromFile replaces the current frequency table by parsing frequencies from the given file path.
func (p *Processor) LoadFrequenciesFromFile(path string) error {
	ft, err := LoadFrequenciesFromFile(path)
	if err != nil {
		return err
	}

	p.ft = ft
	return nil
}

// Save the frequency table to the given file path.
func (p *Processor) Save(path string) error {
	//AJ### TODO: Need to do "atomic" save and replace
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to save the frequency table to file %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close %s. %v", path, err)
		}
	}()

	if err := p.ft.Save(f); err != nil {
		return fmt.Errorf("failed to save the frequency table to file %q. %w", path, err)
	}
	return nil
}

// // UpdateTableByParsingLettersFromFiles updates the frequency table by parsing letter ngrams from the given input paths.
func (p *Processor) ProcessFiles(ctx context.Context, paths []string) error {
	total := len(paths)

	var fn processFileFunc
	if p.mode == ProcessWords {
		fn = func(ctx context.Context, r io.Reader) error {
			return p.ft.ParseWordTokens(ctx, r, p.language, p.tokenSize)
		}
	} else {
		fn = func(ctx context.Context, r io.Reader) error {
			return p.ft.ParseLetterTokens(ctx, r, p.language, p.tokenSize)
		}
	}

	if !p.isNullProgressReporter() {
		totalSize, err := sumFilesizes(paths)
		if err != nil {
			return fmt.Errorf("failed to get the total file size. %w", err)
		}
		p.progressReporter.AddToTotalSize(int64(totalSize))
	}

	for i, path := range paths {
		p.progressReporter.Started(path, i, total)

		if err := p.processFile(ctx, path, fn); err != nil {
			return fmt.Errorf("failed to process the file %q. %w", path, err)
		}
	}

	p.ft.Update()
	return nil
}

//-----------------------------------------------------------------------------

type processFileFunc func(ctx context.Context, r io.Reader) error

func (p *Processor) processFile(ctx context.Context, path string, fn processFileFunc) error {
	// Check if this is a zip file
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".zip" {
		return p.processZipFile(ctx, path, fn)
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

	r := p.progressReporter.Reader(bufio.NewReader(f))
	err = fn(ctx, r)
	if err != nil {
		return err
	}
	return nil
}

func (p *Processor) processZipFile(ctx context.Context, path string, fn processFileFunc) error {
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
	if !p.isNullProgressReporter() {
		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to get the file size for %q. %w", path, err)
		}
		p.progressReporter.AddToTotalSize(-fi.Size())

		totalUncompressedSize := uint64(0)
		for _, f := range zf.File {
			if !filter(f) {
				continue
			}
			totalUncompressedSize += f.UncompressedSize64
		}
		p.progressReporter.AddToTotalSize(int64(totalUncompressedSize))
	}

	// Process each file in the zip
	for _, f := range zf.File {
		if !filter(f) {
			continue
		}

		zr, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %q inside of zip file %q. %w", f.Name, path, err)
		}

		r := p.progressReporter.Reader(bufio.NewReader(zr))
		err = fn(ctx, r)
		if err != nil {
			closer(zr)
			return err
		}
		closer(zr)
	}

	return nil
}

func (p *Processor) isNullProgressReporter() bool {
	_, ok := p.progressReporter.(*nullProgressReporter)
	return ok
}

//-----------------------------------------------------------------------------

func sumFilesizes(paths []string) (uint64, error) {
	total := uint64(0)
	for _, path := range paths {
		fi, err := os.Stat(path)
		if err != nil {
			return 0, fmt.Errorf("failed to get the file size for %q. %w", path, err)
		}
		total += uint64(fi.Size())
	}

	return total, nil
}
