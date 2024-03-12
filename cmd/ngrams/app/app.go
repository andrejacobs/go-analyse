package app

import (
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
	for _, apply := range opts {
		err := apply(&opt)
		if err != nil {
			return nil, fmt.Errorf("failed to configure the app. %w", err)
		}
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
