package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/andrejacobs/go-analyse/internal/compiledinfo"
	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-analyse/text/ngrams"
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

func (a *App) Run(out io.Writer) error {
	if a.opt.update {
		return fmt.Errorf("TODO: implement --update")
	}

	if a.opt.discover {
		return fmt.Errorf("TODO: implement --discover")
	}

	lang, err := a.opt.languages.Get(a.opt.langCode)
	if err != nil {
		return err
	}

	ctx := context.Background()
	var ft *ngrams.FrequencyTable
	if a.opt.words {
		return fmt.Errorf("TODO: implement --words")
	} else {
		ft, err = ngrams.FrequencyTableByParsingLetters(ctx, a.opt.inputs, lang, a.opt.tokenSize)
		if err != nil {
			return err
		}
	}

	if ft != nil {
		if err := a.saveFrequencyTable(ft); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) saveFrequencyTable(ft *ngrams.FrequencyTable) error {
	//AJ### TODO: Need to do "atomic" save and replace (if using update)
	path := a.opt.outPath
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to save the frequency table to file %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close %s. %v", path, err)
		}
	}()

	if err := ft.Save(f); err != nil {
		return fmt.Errorf("failed to save the frequency table to file %q. %w", path, err)
	}
	return nil
}

//-----------------------------------------------------------------------------
// Functional options pattern

type options struct {
	outPath   string
	inputs    []string
	langCode  alphabet.LanguageCode
	languages alphabet.LanguageMap
	words     bool
	tokenSize int
	discover  bool
	update    bool
}

// Option is called to configure the options the app needs to function.
type Option func(opt *options) error

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

// WithDiscoverLanguage configures the app to discover the non-whitespace characters being used.
func WithDiscoverLanguage() Option {
	return func(opt *options) error {
		opt.discover = true
		return nil
	}
}

// WithUpdate configures the app to update an existing ngram output.
func WithUpdate() Option {
	return func(opt *options) error {
		opt.update = true
		return nil
	}
}

// WithOutputPath configures the app to store the ngram output to the given path.
func WithOutputPath(outPath string) Option {
	return func(opt *options) error {
		opt.outPath = outPath
		return nil
	}
}

// WithInputPaths configures the app to parse the given input path files.
func WithInputPaths(paths []string) Option {
	inputs := make([]string, 0, len(paths))
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path != "" {
			inputs = append(inputs, path)
		}
	}

	return func(opt *options) error {
		opt.inputs = inputs
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

	var outPath string
	//TODO: Document that this depends on --discover and --update
	flag.StringVar(&outPath, "o", "", "Path to where the output will be stored.")
	flag.StringVar(&outPath, "out", "", "Path to where the output will be stored.")

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

	var discover bool
	flag.BoolVar(&discover, "d", false, "Discover the non-whitespace letters used and write a languages file to the out path.")
	flag.BoolVar(&discover, "discover", false, "Discover the non-whitespace letters used and write a languages file to the out path.")

	var update bool
	flag.BoolVar(&update, "u", false, "Update the existing ngram output file.")
	flag.BoolVar(&update, "update", false, "Update the existing ngram output file.")

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "Display version information.")
	flag.BoolVar(&showVersion, "version", false, "Display version information.")

	flag.Parse()

	if showVersion {
		printVersion(os.Stdout)
		return nil, ErrExitWithNoErr
	}

	opts = append(opts, WithDefaults())
	opts = append(opts, WithInputPaths(flag.Args()))
	opts = append(opts, WithSize(tokenSize))
	opts = append(opts, WithLanguageCode(alphabet.LanguageCode(langCode)))

	if outPath != "" {
		opts = append(opts, WithOutputPath(outPath))
	}

	if langPath != "" {
		opts = append(opts, WithLanguagesFile(langPath))
	}

	if useLetters {
		opts = append(opts, WithLetters())
	}

	if useWords {
		opts = append(opts, WithWords())
	}

	if discover {
		opts = append(opts, WithDiscoverLanguage())
	}

	if update {
		opts = append(opts, WithUpdate())
	}

	opts = append(opts, resolve())
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

func resolve() Option {
	return func(opt *options) error {
		// ensure at least one input path is given
		if len(opt.inputs) < 1 {
			return fmt.Errorf("expected at least one input path")
		}

		// validate language exists
		if _, ok := opt.languages[opt.langCode]; !ok {
			return fmt.Errorf("failed to find the language %q", opt.langCode)
		}

		// default output path
		if opt.outPath == "" {
			if opt.discover {
				opt.outPath = "./languages.csv"
			} else {
				var ngramType string
				if opt.words {
					ngramType = "words"
				} else {
					ngramType = "letters"
				}
				opt.outPath = fmt.Sprintf("./%s-%s-%d.csv", opt.langCode, ngramType, opt.tokenSize)
			}
		}

		return nil
	}
}

func printVersion(w io.Writer) {
	io.WriteString(w, compiledinfo.UsageNameAndVersion()+"\n")
}
