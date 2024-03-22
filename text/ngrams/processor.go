package ngrams

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/andrejacobs/go-analyse/internal/processor"
	"github.com/andrejacobs/go-analyse/text/alphabet"
)

// FrequencyProcessor is used to parse letter or word ngrams from input sources.
type FrequencyProcessor struct {
	proc      *processor.Processor
	ft        *FrequencyTable
	language  alphabet.Language
	tokenSize int
	mode      ProcessorMode
}

// ProcessorMode specifies whether the processor works on letter or word ngrams.
type ProcessorMode bool

const (
	ProcessLetters ProcessorMode = false
	ProcessWords   ProcessorMode = true
)

// NewFrequencyProcessor creates a new frequency table and does not report progress.
func NewFrequencyProcessor(mode ProcessorMode, language alphabet.Language, tokenSize int) *FrequencyProcessor {
	p := &FrequencyProcessor{
		proc:      processor.NewProcessor(),
		ft:        NewFrequencyTable(),
		language:  language,
		tokenSize: tokenSize,
		mode:      mode,
	}
	return p
}

// SetProgressReporter sets the progress reporter to use.
func (p *FrequencyProcessor) SetProgressReporter(reporter processor.ProgressReporter) {
	p.proc.SetProgressReporter(reporter)
}

// Table returns the frequency table.
func (p *FrequencyProcessor) FrequencyTable() *FrequencyTable {
	return p.ft
}

// LoadFrequenciesFromFile replaces the current frequency table by parsing frequencies from the given file path.
func (p *FrequencyProcessor) LoadFrequenciesFromFile(path string) error {
	ft, err := LoadFrequenciesFromFile(path)
	if err != nil {
		return err
	}

	p.ft = ft
	return nil
}

// Save the frequency table to the given file path.
func (p *FrequencyProcessor) Save(path string) error {
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

// ProcessFiles updates the frequency table by parsing letter or word ngrams from the given input paths.
func (p *FrequencyProcessor) ProcessFiles(ctx context.Context, paths []string) error {
	var fn processor.ProcessFunc
	if p.mode == ProcessWords {
		fn = func(ctx context.Context, r io.Reader) error {
			return p.ft.ParseWordTokens(ctx, r, p.language, p.tokenSize)
		}
	} else {
		fn = func(ctx context.Context, r io.Reader) error {
			return p.ft.ParseLetterTokens(ctx, r, p.language, p.tokenSize)
		}
	}

	if err := p.proc.ProcessFiles(ctx, paths, fn); err != nil {
		return err
	}

	p.ft.Update()
	return nil
}
