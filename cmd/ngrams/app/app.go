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
	"github.com/dustin/go-humanize"
	"github.com/schollz/progressbar/v3"
)

// Main parses the command line arguments and runs the app.
// Decoupled for unit-testing.
func Main(stdOut io.Writer, stdErr io.Writer) error {
	opts, err := parseArgs()
	if err != nil {
		if errors.Is(err, ErrExitWithNoErr) {
			return nil
		}
		fmt.Fprintf(stdErr, "ERROR: %v\n", err)
		return err
	}

	a, err := newApp(opts...)
	if err != nil {
		fmt.Fprintf(stdErr, "ERROR: %v\n", err)
		return err
	}

	if err := a.run(stdOut, stdErr); err != nil {
		fmt.Fprintf(stdErr, "ERROR: %v\n", err)
		return err
	}

	return nil
}

// application provides all the functionality for running the ngrams command line app.
type application struct {
	opt      options
	stdOut   io.Writer
	stdErr   io.Writer
	progress *progressReporter
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

func (a *application) run(stdOut io.Writer, stdErr io.Writer) error {
	ctx := context.Background()

	a.stdOut = stdOut
	a.stdErr = stdErr

	if a.opt.progress {
		a.verbose("Calculating file sizes...\n")
		totalSize, err := sumFilesizes(a.opt.inputs)
		a.verbose("Total size: %s\n", humanize.Bytes(totalSize))
		if err != nil {
			return err
		}

		a.progress = &progressReporter{
			out:       a.stdOut,
			totalSize: totalSize,
		}
	}

	if a.opt.discover {
		return a.discoverLetters(ctx)
	}

	return a.generateNgrams(ctx)
}

func (a *application) generateNgrams(ctx context.Context) error {
	var ft *ngrams.FrequencyTable
	var err error

	if a.opt.update {
		exists, err := pathExists(a.opt.outPath)
		if err != nil {
			return err
		}
		if exists {
			a.verbose("Loading existing frequency table: %q\n", a.opt.outPath)
			ft, err = ngrams.LoadFrequenciesFromFile(a.opt.outPath)
			if err != nil {
				return err
			}
		} else {
			ft = ngrams.NewFrequencyTable()
		}
	} else {
		ft = ngrams.NewFrequencyTable()
	}

	if a.progress != nil {
		ft.SetProgressReporter(a.progress)
	} else if a.opt.verbose {
		ft.SetProgressReporter(&verboseReporter{a: a})
	}

	lang, err := a.opt.languages.Get(a.opt.langCode)
	if err != nil {
		return err
	}

	a.verbose("Language: %s - %s\n", lang.Code, lang.Name)

	if a.opt.words {
		a.verbose("Generating %d word ngrams...\n", a.opt.tokenSize)
		err = ft.UpdateTableByParsingWordsFromFiles(ctx, a.opt.inputs, lang, a.opt.tokenSize)
		if err != nil {
			return err
		}
	} else {
		a.verbose("Generating %d letter ngrams...\n", a.opt.tokenSize)
		err = ft.UpdateTableByParsingLettersFromFiles(ctx, a.opt.inputs, lang, a.opt.tokenSize)
		if err != nil {
			return err
		}
	}

	if ft != nil {
		if err := a.saveFrequencyTable(ft); err != nil {
			return err
		}
		a.verbose("Created frequency table at: %q\n", a.opt.outPath)
	}

	return nil
}

func (a *application) discoverLetters(ctx context.Context) error {
	a.verbose("Discovering letters being used...\n")
	f, err := os.Create(a.opt.outPath)
	if err != nil {
		return fmt.Errorf("failed to create the languages file %q. %w", a.opt.outPath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(a.stdErr, "ERROR: %s\n", err)
		}
	}()

	result := collection.NewSet[rune]()
	total := len(a.opt.inputs)

	closer := func(f io.ReadCloser, path string) {
		if err := f.Close(); err != nil {
			fmt.Fprintf(a.stdErr, "ERROR: failed to close %q. %v", path, err)
		}
	}

	for i, path := range a.opt.inputs {
		if a.progress != nil {
			a.progress.Started(path, i, total)
		} else {
			a.verbose("[%d/%d] %s\n", i+1, total, path)
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open %q. %w", path, err)
		}

		var r io.Reader
		if a.progress != nil {
			r = a.progress.Reader(f)
		} else {
			r = f
		}

		runes, err := alphabet.DiscoverLetters(ctx, r)
		if err != nil {
			closer(f, path)
			return fmt.Errorf("failed to discover the letters from %q. %w", path, err)
		}
		result.InsertSlice(runes)
		closer(f, path)
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

	a.verbose("Created language file at: %q\n", a.opt.outPath)
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
			fmt.Fprintf(a.stdErr, "ERROR: failed to close %s. %v", path, err)
		}
	}()

	if err := ft.Save(f); err != nil {
		return fmt.Errorf("failed to save the frequency table to file %q. %w", path, err)
	}
	return nil
}

func (a *application) verbose(format string, args ...any) {
	if a.opt.verbose {
		fmt.Fprintf(a.stdOut, format, args...)
	}
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

	verbose  bool
	progress bool
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

// withVerbose configures the app to write more information out to Stdout.
func withVerbose() optionFunc {
	return func(opt *options) error {
		opt.verbose = true
		return nil
	}
}

// withProgress configures the app to display progress updates on Stdout.
func withProgress() optionFunc {
	return func(opt *options) error {
		opt.progress = true
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

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Display more information on STDOUT.")

	var progress bool
	flag.BoolVar(&progress, "progress", false, "Display progress updates on STDOUT.")

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

	if verbose {
		opts = append(opts, withVerbose())
	}

	if progress {
		opts = append(opts, withProgress())
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

// Check if the path exists
// If the path exists then (true, nil) is returned
// If the path does not exist then (false, nil) is returned
// If an error occurred while trying to check if the path exists then (false, err) is returned
func pathExists(path string) (bool, error) {
	//NOTE: This is copied from my previous fileutils package (replace this with the new planned repo/modules)
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

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

//-----------------------------------------------------------------------------
// Progress reporting

// With a Progress bar
// Implements ngrams.Progress interface
type progressReporter struct {
	out         io.Writer
	totalSize   uint64
	progressBar *progressbar.ProgressBar
}

func (p *progressReporter) Started(path string, index int, total int) {
	if p.progressBar == nil {
		// Lazy initialized to stop writing progress bar before we actually started
		// otherwise this interrupts other STDOUT printing
		p.progressBar = progressbar.DefaultBytes(int64(p.totalSize))
	}
	p.progressBar.Describe(fmt.Sprintf("[%d/%d]", index+1, total))
}

func (p *progressReporter) Reader(r io.Reader) io.Reader {
	pbr := progressbar.NewReader(r, p.progressBar)
	return &pbr
}

// Only used when verbose is enabled and only
// because I wanted to report which file is being worked on
// Implements ngrams.Progress interface
type verboseReporter struct {
	a *application
}

func (v *verboseReporter) Started(path string, index int, total int) {
	v.a.verbose("[%d/%d] %s\n", index+1, total, path)
}

func (v *verboseReporter) Reader(r io.Reader) io.Reader {
	return r
}
