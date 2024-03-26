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
	"github.com/schollz/progressbar/v3"
)

// Main parses the command line arguments and runs the app.
// Decoupled for unit-testing.
func Main(stdOut io.Writer, stdErr io.Writer) error {
	opts, err := parseArgs(stdOut)
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

// newApp creates a new [application] with the given option configuration.
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
		a.progress = &progressReporter{
			out: a.stdOut,
		}
	}

	if a.opt.discover {
		return a.discoverLetters(ctx)
	}

	return a.generateNgrams(ctx)
}

func (a *application) generateNgrams(ctx context.Context) error {

	lang, err := a.opt.languages.Get(a.opt.langCode)
	if err != nil {
		return err
	}
	a.verbose("Language: %s - %s\n", lang.Code, lang.Name)

	p := ngrams.NewFrequencyProcessor(ngrams.ProcessorMode(a.opt.words), lang, a.opt.tokenSize)

	if a.opt.update {
		exists, err := pathExists(a.opt.outPath)
		if err != nil {
			return err
		}
		if exists {
			a.verbose("Loading existing frequency table: %q\n", a.opt.outPath)
			err = p.LoadFrequenciesFromFile(a.opt.outPath)
			if err != nil {
				return err
			}
		}
	}

	if a.opt.words {
		a.verbose("Generating %d word ngrams...\n", a.opt.tokenSize)
	} else {
		a.verbose("Generating %d letter ngrams...\n", a.opt.tokenSize)
	}

	if a.progress != nil {
		a.progress.progressBar = progressbar.DefaultBytes(1)
		p.SetProgressReporter(a.progress)
	} else if a.opt.verbose {
		p.SetProgressReporter(&verboseReporter{a: a})
	}

	if err = p.ProcessFiles(ctx, a.opt.inputs); err != nil {
		return err
	}

	if a.progress != nil {
		_ = a.progress.progressBar.Finish()
	}

	a.verbose("Saving frequency table...\n")
	if err := p.Save(a.opt.outPath); err != nil {
		return err
	}

	a.verbose("Created frequency table at: %q\n", a.opt.outPath)
	return nil
}

func (a *application) discoverLetters(ctx context.Context) error {
	a.verbose("Discovering letters being used...\n")

	p := alphabet.NewDiscoverProcessor()

	if a.progress != nil {
		a.progress.progressBar = progressbar.DefaultBytes(1)
		p.SetProgressReporter(a.progress)
	} else if a.opt.verbose {
		p.SetProgressReporter(&verboseReporter{a: a})
	}

	if err := p.ProcessFiles(ctx, a.opt.inputs); err != nil {
		return err
	}

	if a.progress != nil {
		_ = a.progress.progressBar.Finish()
	}

	if err := p.Save(a.opt.outPath); err != nil {
		return err
	}

	a.verbose("Created language file at: %q\n", a.opt.outPath)
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

// withLanguageCode configures the language to be used.
func withLanguageCode(code alphabet.LanguageCode) optionFunc {
	return func(opt *options) error {
		opt.langCode = code
		return nil
	}
}

// withLanguagesFile will load the languages from the given file path.
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
func parseArgs(stdOut io.Writer) ([]optionFunc, error) {
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
	flag.BoolVar(&showVersion, "version", false, "Display version information.")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Display more information on STDOUT.")
	flag.BoolVar(&verbose, "verbose", false, "Display more information on STDOUT.")

	var progress bool
	flag.BoolVar(&progress, "progress", false, "Display progress updates on STDOUT.")

	flag.Parse()

	if showVersion {
		printVersion(stdOut)
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
	_, _ = io.WriteString(w, compiledinfo.UsageNameAndVersion()+"\n")
}

// Check if the path exists.
// If the path exists then (true, nil) is returned.
// If the path does not exist then (false, nil) is returned.
// If an error occurred while trying to check if the path exists then (false, err) is returned.
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

//-----------------------------------------------------------------------------
// Progress reporting

// With a Progress bar.
// Implements ngrams.Progress interface.
type progressReporter struct {
	out         io.Writer
	totalSize   int64
	progressBar *progressbar.ProgressBar
}

func (p *progressReporter) Started(path string, index int, total int) {
	p.progressBar.Describe(fmt.Sprintf("[%d/%d]", index+1, total))
}

func (p *progressReporter) Reader(r io.Reader) io.Reader {
	pbr := progressbar.NewReader(r, p.progressBar)
	return &pbr
}

func (p *progressReporter) AddToTotalSize(add int64) {
	if add < 0 {
		// The processor will inform us to subtract the zip file size
		// since the real size is dependant on the total uncompressed size
		// of all the zipped files.
		// However the progress bar will stop updating when it thinks we reached the Max64
		// which does happens just before we get a chance to update the new Max64 size
		// so the hack here is to just be one byte short
		add += 1
	}
	p.totalSize += add
	p.progressBar.ChangeMax64(int64(p.totalSize))
}

// Only used when verbose is enabled and only
// because I wanted to report which file is being worked on.
// Implements ngrams.Progress interface.
type verboseReporter struct {
	a *application
}

func (v *verboseReporter) Started(path string, index int, total int) {
	v.a.verbose("[%d/%d] %s\n", index+1, total, path)
}

func (v *verboseReporter) Reader(r io.Reader) io.Reader {
	return r
}

func (v *verboseReporter) AddToTotalSize(add int64) {
}
