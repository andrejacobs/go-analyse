package app

import (
	"flag"
	"fmt"

	"github.com/andrejacobs/go-analyse/text/alphabet"
)

// App provides all the functionality for running the ngrams command line app.
type App struct {
	opt options
}

// New creates a new [App] with the given option configuration
func New(opts ...Option) (*App, error) {
	var opt options
	if err := applyOptions(&opt, opts); err != nil {
		return nil, err
	}

	result := &App{
		opt: opt,
	}
	return result, nil
}

//-----------------------------------------------------------------------------
// Functional options pattern

type options struct {
	langCode  alphabet.LanguageCode
	words     bool
	tokenSize int
}

// Option is called to configure the options the app needs to function.
type Option func(opt *options) error

// WithDefaults return the default configuration options for the app.
func WithDefaults() Option {
	return func(opt *options) error {
		opt.langCode = "en"
		opt.words = false
		opt.tokenSize = 1
		return nil
	}
}

// WithLanguageCode configures the langauge to be used.
func WithLanguageCode(code alphabet.LanguageCode) Option {
	return func(opt *options) error {
		//AJ### TODO: Need to think about how I specify language files, and do validation
		opt.langCode = code
		return nil
	}
}

// WithLetters configures the app to calculate letter combinations. E.g. bigrams st, er, ao, ie.
func WithLetters() Option {
	return func(opt *options) error {
		opt.words = false
		return nil
	}
}

// WithWords configures the app to calculate word combinations. E.g. bigrams she walked, he jumped.
func WithWords() Option {
	return func(opt *options) error {
		opt.words = true
		return nil
	}
}

// WithSize defines how many letters or words form a single ngram.
func WithSize(size int) Option {
	return func(opt *options) error {
		if size < 1 {
			return fmt.Errorf("invalid ngram size %d", size)
		}
		opt.tokenSize = size
		return nil
	}
}

//-----------------------------------------------------------------------------
// Command line parsing

func ParseArgs() ([]Option, error) {
	opts := make([]Option, 0, 10)

	tokenSize := 0
	flag.IntVar(&tokenSize, "s", 1, "Ngram size. The number of letters or words that form a single ngram.")
	flag.IntVar(&tokenSize, "size", 1, "Ngram size. The number of letters or words that form a single ngram.")

	langCode := ""
	flag.StringVar(&langCode, "a", "en", "Alphabet language code. E.g. en = English")
	flag.StringVar(&langCode, "language", "en", "Alphabet language code. E.g. en = English")

	flag.Parse()

	opts = append(opts, WithDefaults())
	opts = append(opts, WithSize(tokenSize))
	opts = append(opts, WithLanguageCode(alphabet.LanguageCode(langCode)))

	return opts, nil
}

//-----------------------------------------------------------------------------

func applyOptions(opt *options, opts []Option) error {
	for _, apply := range opts {
		err := apply(opt)
		if err != nil {
			return fmt.Errorf("failed to configure the app. %w", err)
		}
	}
	return nil
}
