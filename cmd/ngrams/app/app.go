package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/andrejacobs/go-analyse/internal/compiledinfo"
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
	languages alphabet.LanguageMap
	words     bool
	tokenSize int
}

// Option is called to configure the options the app needs to function.
type Option func(opt *options) error

//AJ### TODO: Add a Validate func, append at the end of the options. Check that the langCode is in the map

// WithDefaults return the default configuration options for the app.
func WithDefaults() Option {
	return func(opt *options) error {
		opt.langCode = "en"
		opt.languages = alphabet.BuiltinLanguages()
		opt.words = false
		opt.tokenSize = 1
		return nil
	}
}

// WithLanguageCode configures the langauge to be used.
func WithLanguageCode(code alphabet.LanguageCode) Option {
	return func(opt *options) error {
		opt.langCode = code
		return nil
	}
}

// WithLanguagesFile will load the languages from the given file path
func WithLanguagesFile(path string) Option {
	return func(opt *options) error {
		var err error
		opt.languages, err = alphabet.LoadLanguagesFromFile(path)
		if err != nil {
			return err
		}
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

var ErrExitWithNoErr = errors.New("not an error")

// ParseArgs will parse the command line arguments and create the slice of options required
// to create the app.
func ParseArgs() ([]Option, error) {
	opts := make([]Option, 0, 10)

	var tokenSize int
	flag.IntVar(&tokenSize, "s", 1, "Ngram size. The number of letters or words that form a single ngram.")
	flag.IntVar(&tokenSize, "size", 1, "Ngram size. The number of letters or words that form a single ngram.")

	var langCode string
	flag.StringVar(&langCode, "a", "en", "Alphabet language code. E.g. en = English")
	flag.StringVar(&langCode, "lang", "en", "Alphabet language code. E.g. en = English")

	var langPath string
	//TODO: Document the CSV format: #code,name,letters
	flag.StringVar(&langPath, "languages", "", "Path to a languages definition file.")

	var useLetters bool
	flag.BoolVar(&useLetters, "l", true, "Create letter ngram combinations. E.g. bigrams st,er,ae,ie.")
	flag.BoolVar(&useLetters, "letters", true, "Create letter ngram combinations. E.g. bigrams st,er,ae,ie.")

	var useWords bool
	flag.BoolVar(&useWords, "w", false, "Create word ngram combinations. E.g. bigrams he jumped, she walked")
	flag.BoolVar(&useWords, "words", false, "Create word ngram combinations. E.g. bigrams he jumped, she walked")

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "Display version information.")
	flag.BoolVar(&showVersion, "version", false, "Display version information.")

	flag.Parse()

	if showVersion {
		printVersion(os.Stdout)
		return nil, ErrExitWithNoErr
	}

	opts = append(opts, WithDefaults())
	opts = append(opts, WithSize(tokenSize))
	opts = append(opts, WithLanguageCode(alphabet.LanguageCode(langCode)))

	if useLetters {
		opts = append(opts, WithLetters())
	}

	if useWords {
		opts = append(opts, WithWords())
	}

	if langPath != "" {
		opts = append(opts, WithLanguagesFile(langPath))
	}

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

func printVersion(w io.Writer) {
	io.WriteString(w, compiledinfo.UsageNameAndVersion()+"\n")
}
