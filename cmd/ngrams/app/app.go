package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/andrejacobs/go-analyse/internal/compiledinfo"
	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-analyse/text/ngrams"
	"github.com/andrejacobs/go-collection/collection"
)

// Main parses the command line arguments and runs the app.
// Decoupled for unit-testing.
func Main(out io.Writer, errOut io.Writer) error {
	opts, err := parseArgs()
	if err != nil {
		fmt.Fprintf(errOut, "ERROR: %v\n", err)
		return err
	}

	a, err := newApp(opts...)
	if err != nil {
		fmt.Fprintf(errOut, "ERROR: %v\n", err)
		return err
	}

	if err := a.run(out, errOut); err != nil {
		fmt.Fprintf(errOut, "ERROR: %v\n", err)
		return err
	}

	return nil
}

// application provides all the functionality for running the ngrams command line app.
type application struct {
	opt options
}

// newApp creates a new [application] with the given option configuration
func newApp(opts ...optionFunc) (*application, error) {
	var opt options
	if err := applyOptions(&opt, opts); err != nil {
		return nil, err
	}

	result := &application{
		opt: opt,
	}

	return result, nil
}

func (a *application) run(out io.Writer, errOut io.Writer) error {
	ctx := context.Background()

	if a.opt.discover {
		return a.discoverLetters(ctx, out, errOut)
	}

	return a.generateNgrams(ctx, out, errOut)
}

func (a *application) generateNgrams(ctx context.Context, out io.Writer, errOut io.Writer) error {
	if a.opt.update {
		return fmt.Errorf("TODO: implement --update")
	}

	lang, err := a.opt.languages.Get(a.opt.langCode)
	if err != nil {
		return err
	}

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

func (a *application) discoverLetters(ctx context.Context, out io.Writer, errOut io.Writer) error {
	f, err := os.Create(a.opt.outPath)
	if err != nil {
		return fmt.Errorf("failed to create the languages file %q. %w", a.opt.outPath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}
	}()

	result := collection.NewSet[rune]()

	for _, path := range a.opt.inputs {
		runes, err := alphabet.DiscoverLettersFromFile(ctx, path)
		if err != nil {
			return fmt.Errorf("failed to discover the letters from %q. %w", path, err)
		}
		result.InsertSlice(runes)
	}

	_, err = io.WriteString(f, "#code,name,letters\n")
	if err != nil {
		return fmt.Errorf("failed to write csv header to %q. %w", a.opt.outPath, err)
	}

	runes := result.Items()
	slices.Sort(runes)
	letters := strings.ReplaceAll(string(runes), `"`, `""`)

	_, err = io.WriteString(f, `unknown,unknown,"`+letters+`"`)
	if err != nil {
		return fmt.Errorf("failed to write csv header to %q. %w", a.opt.outPath, err)
	}

	return nil
}

func (a *application) saveFrequencyTable(ft *ngrams.FrequencyTable) error {
	//AJ### TODO: Need to do "atomic" save and replace (if using update)
	path := a.opt.outPath
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to save the frequency table to file %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close %s. %v", path, err)
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

// optionFunc is called to configure the options the app needs to function.
type optionFunc func(opt *options) error

// withDefaults return the default configuration options for the app.
func withDefaults() optionFunc {
	return func(opt *options) error {
		opt.langCode = "en"
		opt.languages = alphabet.BuiltinLanguages()
		opt.words = false
		opt.tokenSize = 1
		return nil
	}
}

// withLanguageCode configures the langauge to be used.
func withLanguageCode(code alphabet.LanguageCode) optionFunc {
	return func(opt *options) error {
		opt.langCode = code
		return nil
	}
}

// withLanguagesFile will load the languages from the given file path
func withLanguagesFile(path string) optionFunc {
	return func(opt *options) error {
		var err error
		opt.languages, err = alphabet.LoadLanguagesFromFile(path)
		if err != nil {
			return err
		}
		return nil
	}
}

// withLetters configures the app to calculate letter combinations. E.g. bigrams st, er, ao, ie.
func withLetters() optionFunc {
	return func(opt *options) error {
		opt.words = false
		return nil
	}
}

// withWords configures the app to calculate word combinations. E.g. bigrams she walked, he jumped.
func withWords() optionFunc {
	return func(opt *options) error {
		opt.words = true
		return nil
	}
}

// withSize defines how many letters or words form a single ngram.
func withSize(size int) optionFunc {
	return func(opt *options) error {
		if size < 1 {
			return fmt.Errorf("invalid ngram size %d", size)
		}
		opt.tokenSize = size
		return nil
	}
}

// withDiscoverLanguage configures the app to discover the non-whitespace characters being used.
func withDiscoverLanguage() optionFunc {
	return func(opt *options) error {
		opt.discover = true
		return nil
	}
}

// withUpdate configures the app to update an existing ngram output.
func withUpdate() optionFunc {
	return func(opt *options) error {
		opt.update = true
		return nil
	}
}

// withOutputPath configures the app to store the ngram output to the given path.
func withOutputPath(outPath string) optionFunc {
	return func(opt *options) error {
		opt.outPath = outPath
		return nil
	}
}

// withInputPaths configures the app to parse the given input path files.
func withInputPaths(paths []string) optionFunc {
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

// parseArgs will parse the command line arguments and create the slice of options required
// to create the app.
func parseArgs() ([]optionFunc, error) {
	opts := make([]optionFunc, 0, 10)

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

	opts = append(opts, withDefaults())
	opts = append(opts, withInputPaths(flag.Args()))
	opts = append(opts, withSize(tokenSize))
	opts = append(opts, withLanguageCode(alphabet.LanguageCode(langCode)))

	if outPath != "" {
		opts = append(opts, withOutputPath(outPath))
	}

	if langPath != "" {
		opts = append(opts, withLanguagesFile(langPath))
	}

	if useLetters {
		opts = append(opts, withLetters())
	}

	if useWords {
		opts = append(opts, withWords())
	}

	if discover {
		opts = append(opts, withDiscoverLanguage())
	}

	if update {
		opts = append(opts, withUpdate())
	}

	opts = append(opts, resolve())
	return opts, nil
}

//-----------------------------------------------------------------------------

func applyOptions(opt *options, opts []optionFunc) error {
	for _, apply := range opts {
		err := apply(opt)
		if err != nil {
			return fmt.Errorf("failed to configure the app. %w", err)
		}
	}
	return nil
}

func resolve() optionFunc {
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
