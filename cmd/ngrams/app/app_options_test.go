package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionsWithDefaults(t *testing.T) {
	var opt options
	require.NoError(t, withDefaults()(&opt))

	assert.Equal(t, alphabet.LanguageCode("en"), opt.langCode)
	assert.Equal(t, alphabet.BuiltinLanguages(), opt.languages)
	assert.Equal(t, 1, opt.tokenSize)
	assert.False(t, opt.words)
	assert.False(t, opt.discover)
	assert.False(t, opt.update)
	assert.Equal(t, "", opt.outPath)
	assert.Empty(t, opt.inputs)
	assert.False(t, opt.verbose)
	assert.False(t, opt.progress)
}

func TestParseArgs(t *testing.T) {
	backupArgs := os.Args
	defer func() {
		os.Args = backupArgs
	}()

	invalidLanguages := invalidLanguagesFile(t)
	defer os.Remove(invalidLanguages)

	validLanguages := validLanguagesFile(t)
	defer os.Remove(validLanguages)

	testCases := []struct {
		desc       string
		args       string
		expected   []optionFunc
		errMsg     string
		assertFunc func(t *testing.T, opt *options)
	}{
		{desc: "invalid size: -s", args: "-s 0 ./in.txt", errMsg: "invalid ngram size 0"},
		{desc: "invalid size: --size", args: "--size 0 ./in.txt", errMsg: "invalid ngram size 0"},
		{desc: "size: -s 4", args: "-s 4 ./in.txt", expected: []optionFunc{withSize(4)}},
		{desc: "size: --size 4", args: "--size 4 ./in.txt", expected: []optionFunc{withSize(4)}},

		{desc: "language: --lang af", args: "--lang af ./in.txt", expected: []optionFunc{withLanguageCode("af")}},
		{desc: "language: -a en", args: "-a en ./in.txt", expected: []optionFunc{withLanguageCode("en")}},

		{desc: "invalid languages file: --languages",
			args:   fmt.Sprintf("--languages %s ./in.txt", invalidLanguages),
			errMsg: "failed to configure the app. failed to load languages from"},

		{desc: "valid languages file: --languages",
			args:     fmt.Sprintf("--languages %s ./in.txt", validLanguages),
			expected: []optionFunc{withLanguagesFile(validLanguages)},
			assertFunc: func(t *testing.T, opt *options) {
				_, ok := opt.languages["coding"]
				assert.True(t, ok)
			},
		},

		{desc: "missing language: --languages -a af",
			args:   fmt.Sprintf("--languages %s -a af ./in.txt", validLanguages),
			errMsg: "failed to configure the app. failed to find the language \"af\"",
		},

		{desc: "letters: -l", args: "-l ./in.txt", expected: []optionFunc{withLetters()}},
		{desc: "letters: --letters", args: "--letters ./in.txt", expected: []optionFunc{withLetters()}},
		{desc: "words: -w", args: "-w ./in.txt", expected: []optionFunc{withWords()}},
		{desc: "words: --words", args: "--words ./in.txt", expected: []optionFunc{withWords()}},
		{desc: "mixing letters and words: -w -l", args: "-w -l ./in.txt", expected: []optionFunc{withWords()}},
		{desc: "mixing letters and words: -l -w", args: "-l -w ./in.txt", expected: []optionFunc{withWords()}},

		{desc: "discover: -d", args: "-d ./in.txt", expected: []optionFunc{withDiscoverLanguage()}},
		{desc: "discover: --discover", args: "--discover ./in.txt", expected: []optionFunc{withDiscoverLanguage()}},

		{desc: "update: -u", args: "-u ./in.txt", expected: []optionFunc{withUpdate()}},
		{desc: "update: --update", args: "--update ./in.txt", expected: []optionFunc{withUpdate()}},

		{desc: "output: -o", args: "-o ./test.csv ./in.txt", expected: []optionFunc{withOutputPath("./test.csv")}},
		{desc: "output: --out", args: "--out ./test.csv ./in.txt", expected: []optionFunc{withOutputPath("./test.csv")}},

		{desc: "default output", args: "./in.txt", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, "./en-letters-1.csv", opt.outPath)
		}},
		{desc: "default output: -s 2", args: "-s 2 ./in.txt", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, "./en-letters-2.csv", opt.outPath)
		}},
		{desc: "default output: -s 3 -w -a af", args: "-s 3 -w -a af ./in.txt", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, "./af-words-3.csv", opt.outPath)
		}},
		{desc: "default output: -d", args: "-d ./in.txt", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, "./languages.csv", opt.outPath)
		}},

		{desc: "no input paths", args: "", errMsg: "failed to configure the app. expected at least one input path"},
		{desc: "input path", args: "./input1.txt ./input2.txt", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, []string{"./input1.txt", "./input2.txt"}, opt.inputs)
		}},

		{desc: "verbose: --verbose", args: "--verbose ./in.txt", expected: []optionFunc{withVerbose()}},
		{desc: "progress: --progress", args: "--progress ./in.txt", expected: []optionFunc{withProgress()}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Fake CLI args for flag package
			os.Args = make([]string, 0)
			os.Args = append(os.Args, "ngrams")
			os.Args = append(os.Args, strings.Split(tC.args, " ")...)

			// Parse the args
			opts, err := parseArgs()
			require.NoError(t, err)
			// Reset flag package for next test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Apply the options (as the app would)
			var opt options
			err = applyOptions(&opt, opts)

			// Check for expected error
			if tC.errMsg != "" {
				assert.ErrorContains(t, err, tC.errMsg)
			}

			// Check for expected options to be applied
			if len(tC.expected) > 0 {
				var expectedOpt options
				require.NoError(t, applyOptions(&expectedOpt, []optionFunc{withDefaults()}))
				require.NoError(t, applyOptions(&expectedOpt, tC.expected))
				require.NoError(t, applyOptions(&expectedOpt, []optionFunc{withInputPaths([]string{"./in.txt"})}))
				require.NoError(t, applyOptions(&expectedOpt, []optionFunc{resolve()}))

				assert.Equal(t, expectedOpt, opt)
			}

			// Perform custom assert checks
			if tC.assertFunc != nil {
				tC.assertFunc(t, &opt)
			}
		})
	}
}

func invalidLanguagesFile(t *testing.T) string {
	f, err := os.CreateTemp("", "invalid-lang.csv")
	require.NoError(t, err)
	defer f.Close()
	return f.Name()
}

func validLanguagesFile(t *testing.T) string {
	f, err := os.CreateTemp("", "valid-lang.csv")
	require.NoError(t, err)
	defer f.Close()

	_, _ = io.WriteString(f, `#code,name,letters
en,English,abcdefghijklmnopqrstuvwxyz
coding,Coding,{}[]()/$
`)

	return f.Name()
}
