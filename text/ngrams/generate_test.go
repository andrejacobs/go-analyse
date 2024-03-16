//go:build gendata

// Generate the testdata
// Run: go test -v -tags=gendata -run "TestGenerate.+" github.com/andrejacobs/go-analyse/text/ngrams

package ngrams_test

import (
	"context"
	"os"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-analyse/text/ngrams"
	"github.com/stretchr/testify/require"
)

// Use this to generate the test data
// Not ideal to be doing chicken-n-egg and generating the testdata using the code
// you actually want to test. I will just have to be 100% sure of the results
func TestGenerateLetterFrequencies(t *testing.T) {
	genFn := func(input string, output string, langCode alphabet.LanguageCode, tokenSize int) {
		lang, err := alphabet.Builtin(langCode)
		require.NoError(t, err)
		ft := ngrams.NewFrequencyTable()
		err = ft.UpdateTableByParsingLettersFromFiles(context.Background(),
			[]string{input}, lang, tokenSize)
		require.NoError(t, err)
		out, err := os.Create(output)
		require.NoError(t, err)
		ft.Save(out)
		out.Close()
	}

	genFn("testdata/en-control.txt", "testdata/freq-1-en-control.csv", "en", 1)
	genFn("testdata/en-control.txt", "testdata/freq-2-en-control.csv", "en", 2)
	genFn("testdata/af-control.txt", "testdata/freq-1-af-control.csv", "af", 1)
	genFn("testdata/af-control.txt", "testdata/freq-2-af-control.csv", "af", 2)

	genFn("testdata/en-alice-partial.txt", "testdata/freq-1-en-alice.csv", "en", 1)
	genFn("testdata/en-alice-partial.txt", "testdata/freq-2-en-alice.csv", "en", 2)
	genFn("testdata/en-alice-partial.txt", "testdata/freq-3-en-alice.csv", "en", 3)

	genFn("testdata/fr-alice-partial.txt", "testdata/freq-1-fr-alice.csv", "fr", 1)
	genFn("testdata/fr-alice-partial.txt", "testdata/freq-2-fr-alice.csv", "fr", 2)
	genFn("testdata/fr-alice-partial.txt", "testdata/freq-3-fr-alice.csv", "fr", 3)
}

// Use this to generate the test data
// Not ideal to be doing chicken-n-egg and generating the testdata using the code
// you actually want to test. I will just have to be 100% sure of the results
func TestGenerateWordFrequencies(t *testing.T) {
	genFn := func(input string, output string, langCode alphabet.LanguageCode, tokenSize int) {
		lang, err := alphabet.Builtin(langCode)
		require.NoError(t, err)
		ft := ngrams.NewFrequencyTable()
		err = ft.UpdateTableByParsingWordsFromFiles(context.Background(),
			[]string{input}, lang, tokenSize)
		require.NoError(t, err)
		out, err := os.Create(output)
		require.NoError(t, err)
		ft.Save(out)
		out.Close()
	}

	genFn("testdata/en-control.txt", "testdata/freq-1w-en-control.csv", "en", 1)
	genFn("testdata/en-control.txt", "testdata/freq-2w-en-control.csv", "en", 2)
	genFn("testdata/af-control.txt", "testdata/freq-1w-af-control.csv", "af", 1)
	genFn("testdata/af-control.txt", "testdata/freq-2w-af-control.csv", "af", 2)

	genFn("testdata/en-alice-partial.txt", "testdata/freq-1w-en-alice.csv", "en", 1)
	genFn("testdata/en-alice-partial.txt", "testdata/freq-2w-en-alice.csv", "en", 2)
	genFn("testdata/en-alice-partial.txt", "testdata/freq-3w-en-alice.csv", "en", 3)

	genFn("testdata/fr-alice-partial.txt", "testdata/freq-1w-fr-alice.csv", "fr", 1)
	genFn("testdata/fr-alice-partial.txt", "testdata/freq-2w-fr-alice.csv", "fr", 2)
	genFn("testdata/fr-alice-partial.txt", "testdata/freq-3w-fr-alice.csv", "fr", 3)
}
