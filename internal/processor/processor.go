// Package processor provides a way to run a function against a collection of various input sources.
package processor

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ProcessFunc is provided to the processor and will be called on each input source that needs processing.
type ProcessFunc func(ctx context.Context, r io.Reader) error

// Processor is used to process multiple input sources.
type Processor struct {
	progressReporter ProgressReporter
}

// NewProcessor creates a new processor.
func NewProcessor() *Processor {
	p := &Processor{
		progressReporter: &nullProgressReporter{},
	}
	return p
}

// SetProgressReporter sets the progress reporter to use.
func (p *Processor) SetProgressReporter(reporter ProgressReporter) {
	p.progressReporter = reporter
}

// ProcessFiles will run the given [ProcessFunc] on the set of input file paths.
// Zip files are also supported and each individual file in the zip will be processed.
func (p *Processor) ProcessFiles(ctx context.Context, paths []string, fn ProcessFunc) error {
	total := len(paths)

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

	return nil
}

//-----------------------------------------------------------------------------

func (p *Processor) processFile(ctx context.Context, path string, fn ProcessFunc) error {
	// Check if this is a zip file
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".zip" {
		return p.processZipFile(ctx, path, fn)
	}
	//AJ### TODO: Add tar.gz support?

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

func (p *Processor) processZipFile(ctx context.Context, path string, fn ProcessFunc) error {
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
