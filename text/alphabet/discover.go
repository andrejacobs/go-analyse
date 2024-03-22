// Copyright (c) 2024 Andre Jacobs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package alphabet

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"unicode"

	"github.com/andrejacobs/go-analyse/internal/processor"
	"github.com/andrejacobs/go-collection/collection"
)

// DiscoverLetters produces a slice containing the unique non-whitespace lowercased letters found in the io.Reader.
func DiscoverLetters(ctx context.Context, input io.Reader) ([]rune, error) {

	result := collection.NewSet[rune]()
	rd := bufio.NewReader(input)

loop:
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return result.Items(), err
			}
			break loop
		default:
			r, _, err := rd.ReadRune()
			if err != nil {
				if err == io.EOF {
					// Done reading
					break loop
				}
				return result.Items(), err
			}

			// Ignore white space
			if unicode.IsSpace(r) {
				continue
			}

			result.Insert(unicode.ToLower(r))
		}
	}

	return result.Items(), nil
}

// DiscoverLettersFromFile produces a slice containing the unique non-whitespace lowercased letters found in the file.
func DiscoverLettersFromFile(ctx context.Context, path string) ([]rune, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close %q. %v", path, err)
		}
	}()

	result, err := DiscoverLetters(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to discover letters from %q. %w", path, err)
	}

	return result, nil
}

//-----------------------------------------------------------------------------

// DiscoverProcessor is used to discover the unique non-whitespace lowercased letters found in the input sources.
type DiscoverProcessor struct {
	proc    *processor.Processor
	letters collection.Set[rune]
}

// NewDiscoverProcessor creates a new processor and does not report progress.
func NewDiscoverProcessor() *DiscoverProcessor {
	p := &DiscoverProcessor{
		proc:    processor.NewProcessor(),
		letters: collection.NewSet[rune](),
	}
	return p
}

// SetProgressReporter sets the progress reporter to use.
func (p *DiscoverProcessor) SetProgressReporter(reporter processor.ProgressReporter) {
	p.proc.SetProgressReporter(reporter)
}

// Letters return the discovered runes. Sounds like a tomb raider story :-D.
func (p *DiscoverProcessor) Letters() []rune {
	return p.letters.Items()
}

// ProcessFiles updates the discovered letters from the given input paths.
func (p *DiscoverProcessor) ProcessFiles(ctx context.Context, paths []string) error {
	err := p.proc.ProcessFiles(ctx, paths, func(ctx context.Context, r io.Reader) error {
		runes, err := DiscoverLetters(ctx, r)
		if err != nil {
			return err
		}
		p.letters.InsertSlice(runes)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Save the languages file to the given file path.
func (p *DiscoverProcessor) Save(path string) error {
	//AJ### TODO: Need to do "atomic" save and replace
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to save the languages file %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close %s. %v", path, err)
		}
	}()

	_, err = io.WriteString(f, "#code,name,letters\n")
	if err != nil {
		return fmt.Errorf("failed to write csv header to %q. %w", path, err)
	}

	runes := p.Letters()
	slices.Sort(runes)
	letters := strings.ReplaceAll(string(runes), `"`, `""`)

	_, err = io.WriteString(f, `unknown,unknown,"`+letters+`"`)
	if err != nil {
		return fmt.Errorf("failed to write csv header to %q. %w", path, err)
	}

	return nil
}
