package app

import (
	"flag"
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
	assert.False(t, opt.words)
	assert.Equal(t, 1, opt.tokenSize)
}

func TestOptionsWithLettersOrWords(t *testing.T) {
	var opt options
	require.NoError(t, WithLetters()(&opt))
	assert.False(t, opt.words)

	require.NoError(t, WithWords()(&opt))
	assert.True(t, opt.words)
}

func TestParseArgs(t *testing.T) {
	backupArgs := os.Args
	defer func() {
		os.Args = backupArgs
	}()

	testCases := []struct {
		desc     string
		args     string
		expected []Option
		errMsg   string
	}{
		{desc: "invalid size: -s", args: "-s 0", errMsg: "invalid ngram size 0"},
		{desc: "invalid size: --size", args: "--size 0", errMsg: "invalid ngram size 0"},
		{desc: "size: -s 4", args: "-s 4", expected: []Option{WithSize(4)}},
		{desc: "size: --size 4", args: "--size 4", expected: []Option{WithSize(4)}},

		{desc: "language: --language af", args: "--language af", expected: []Option{WithLanguageCode("af")}},
		{desc: "language: -a en", args: "-a en", expected: []Option{WithLanguageCode("en")}},

		{desc: "letters: -l", args: "-l", expected: []Option{WithLetters()}},
		{desc: "letters: --letters", args: "--letters", expected: []Option{WithLetters()}},
		{desc: "words: -w", args: "-w", expected: []Option{WithWords()}},
		{desc: "words: --words", args: "--words", expected: []Option{WithWords()}},
		{desc: "mixing letters and words: -w -l", args: "-w -l", expected: []Option{WithWords()}},
		{desc: "mixing letters and words: -l -w", args: "-l -w", expected: []Option{WithWords()}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Fake CLI args for flag package
			os.Args = make([]string, 0)
			os.Args = append(os.Args, "ngrams")
			os.Args = append(os.Args, strings.Split(tC.args, " ")...)

			opts, err := ParseArgs()
			require.NoError(t, err)
			// Reset flag package for next test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			var opt options
			err = applyOptions(&opt, opts)
			if tC.errMsg != "" {
				assert.ErrorContains(t, err, tC.errMsg)
			}

			var expectedOpt options
			require.NoError(t, applyOptions(&expectedOpt, []Option{WithDefaults()}))
			require.NoError(t, applyOptions(&expectedOpt, tC.expected))

			assert.Equal(t, expectedOpt, opt)
		})
	}
}
