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
	require.NoError(t, WithDefaults()(&opt))

	assert.Equal(t, alphabet.LanguageCode("en"), opt.langCode)
	assert.Equal(t, alphabet.BuiltinLanguages(), opt.languages)
	assert.Equal(t, 1, opt.tokenSize)
	assert.False(t, opt.words)
	assert.False(t, opt.discover)
	assert.False(t, opt.update)
	assert.Equal(t, "", opt.outPath)
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
		expected   []Option
		errMsg     string
		assertFunc func(t *testing.T, opt *options)
	}{
		{desc: "invalid size: -s", args: "-s 0", errMsg: "invalid ngram size 0"},
		{desc: "invalid size: --size", args: "--size 0", errMsg: "invalid ngram size 0"},
		{desc: "size: -s 4", args: "-s 4", expected: []Option{WithSize(4)}},
		{desc: "size: --size 4", args: "--size 4", expected: []Option{WithSize(4)}},

		{desc: "language: --lang af", args: "--lang af", expected: []Option{WithLanguageCode("af")}},
		{desc: "language: -a en", args: "-a en", expected: []Option{WithLanguageCode("en")}},

		{desc: "invalid languages file: --languages",
			args:   fmt.Sprintf("--languages %s", invalidLanguages),
			errMsg: "failed to configure the app. failed to load languages from"},

		{desc: "valid languages file: --languages",
			args:     fmt.Sprintf("--languages %s", validLanguages),
			expected: []Option{WithLanguagesFile(validLanguages)},
			assertFunc: func(t *testing.T, opt *options) {
				_, ok := opt.languages["coding"]
				assert.True(t, ok)
			},
		},

		{desc: "missing language: --languages -a af",
			args:   fmt.Sprintf("--languages %s -a af", validLanguages),
			errMsg: "failed to configure the app. failed to find the language \"af\"",
		},

		{desc: "letters: -l", args: "-l", expected: []Option{WithLetters()}},
		{desc: "letters: --letters", args: "--letters", expected: []Option{WithLetters()}},
		{desc: "words: -w", args: "-w", expected: []Option{WithWords()}},
		{desc: "words: --words", args: "--words", expected: []Option{WithWords()}},
		{desc: "mixing letters and words: -w -l", args: "-w -l", expected: []Option{WithWords()}},
		{desc: "mixing letters and words: -l -w", args: "-l -w", expected: []Option{WithWords()}},

		{desc: "discover: -d", args: "-d", expected: []Option{WithDiscoverLanguage()}},
		{desc: "discover: --discover", args: "--discover", expected: []Option{WithDiscoverLanguage()}},

		{desc: "update: -u", args: "-u", expected: []Option{WithUpdate()}},
		{desc: "update: --update", args: "--update", expected: []Option{WithUpdate()}},

		{desc: "output: -o", args: "-o ./test.csv", expected: []Option{WithOutputPath("./test.csv")}},
		{desc: "output: --out", args: "--out ./test.csv", expected: []Option{WithOutputPath("./test.csv")}},

		{desc: "default output", args: "", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, "./en-letters-1.csv", opt.outPath)
		}},
		{desc: "default output: -s 2", args: "-s 2", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, "./en-letters-2.csv", opt.outPath)
		}},
		{desc: "default output: -s 3 -w -a af", args: "-s 3 -w -a af", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, "./af-words-3.csv", opt.outPath)
		}},
		{desc: "default output: -d", args: "-d", assertFunc: func(t *testing.T, opt *options) {
			assert.Equal(t, "./languages.csv", opt.outPath)
		}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Fake CLI args for flag package
			os.Args = make([]string, 0)
			os.Args = append(os.Args, "ngrams")
			os.Args = append(os.Args, strings.Split(tC.args, " ")...)

			// Parse the args
			opts, err := ParseArgs()
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
				require.NoError(t, applyOptions(&expectedOpt, []Option{WithDefaults()}))
				require.NoError(t, applyOptions(&expectedOpt, tC.expected))
				require.NoError(t, applyOptions(&expectedOpt, []Option{resolve()}))

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

	io.WriteString(f, `#code,name,letters
en,English,abcdefghijklmnopqrstuvwxyz
coding,Coding,{}[]()/$
`)

	return f.Name()
}
