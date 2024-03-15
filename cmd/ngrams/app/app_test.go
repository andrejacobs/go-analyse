package app_test

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/cmd/ngrams/app"
	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-analyse/text/ngrams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	backupArgs := os.Args
	defer func() {
		os.Args = backupArgs
	}()

	invalidLanguages := invalidLanguagesFile(t)
	defer os.Remove(invalidLanguages)

	validLanguages := validLanguagesFile(t)
	defer os.Remove(validLanguages)

	outPath := tempOutputPath()
	defer os.Remove(outPath)

	testCases := []struct {
		desc     string
		args     string
		testFunc func(t *testing.T)
	}{
		//------------------------
		// Sad paths

		{desc: "no input args", args: "", testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, "failed to configure the app. expected at least one input path")
			assert.Error(t, err)
		}},

		{desc: "input file does not exist", args: "./in.txt", testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, "open ./in.txt: no such file or directory")
			assert.Error(t, err)
		}},

		//AJ### TODO: an invalid input file (i.e. not parsing)

		{desc: "languages file does not exist", args: "--languages ./lang.txt ./in.txt", testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, "open ./lang.txt: no such file or directory")
			assert.Error(t, err)
		}},

		{desc: "invalid languages file", args: fmt.Sprintf("--languages %s ./in.txt", invalidLanguages), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, fmt.Sprintf("failed to load languages from \"%s\"", invalidLanguages))
			assert.Error(t, err)
		}},

		{desc: "missing built-in language", args: "-a golang ./in.txt", testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, "failed to find the language \"golang\"")
			assert.Error(t, err)
		}},

		{desc: "missing language", args: fmt.Sprintf("-a golang --languages %s ./in.txt", validLanguages), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, "failed to find the language \"golang\"")
			assert.Error(t, err)
		}},

		{desc: "invalid size", args: "-s 0 ./in.txt", testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, "invalid ngram size 0")
			assert.Error(t, err)
		}},

		//------------------------
		// Happy paths

		// Letters

		{desc: "monograms en-control", args: fmt.Sprintf("-s 1 -o %s %s", outPath, inputENControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENControl1)
		}},

		{desc: "bigrams en-control", args: fmt.Sprintf("-s 2 -o %s %s", outPath, inputENControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENControl2)
		}},

		{desc: "monograms af-control", args: fmt.Sprintf("-a af -s 1 -o %s %s", outPath, inputAFControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputAFControl1)
		}},

		{desc: "bigrams af-control", args: fmt.Sprintf("-a af -s 2 -o %s %s", outPath, inputAFControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputAFControl2)
		}},

		{desc: "bigrams en-alice-partial", args: fmt.Sprintf("-s 2 -o %s %s", outPath, inputENAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENAlice2)
		}},

		{desc: "trigrams en-alice-partial", args: fmt.Sprintf("-s 3 -o %s %s", outPath, inputENAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENAlice3)
		}},

		{desc: "bigrams fr-alice-partial", args: fmt.Sprintf("-a fr -s 2 -o %s %s", outPath, inputFRAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputFRAlice2)
		}},

		{desc: "trigrams fr-alice-partial", args: fmt.Sprintf("-a fr -s 3 -o %s %s", outPath, inputFRAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputFRAlice3)
		}},

		// Words

		{desc: "word monograms af-control", args: fmt.Sprintf("-w -a af -s 2 -o %s %s", outPath, inputAFControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputAFControlW2)
		}},

		{desc: "word bigrams en-alice-partial", args: fmt.Sprintf("-w -s 2 -o %s %s", outPath, inputENAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENAliceW2)
		}},

		{desc: "word trigrams en-alice-partial", args: fmt.Sprintf("-w -s 3 -o %s %s", outPath, inputENAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENAliceW3)
		}},

		{desc: "word bigrams fr-alice-partial", args: fmt.Sprintf("-w -a fr -s 2 -o %s %s", outPath, inputFRAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputFRAliceW2)
		}},

		{desc: "word trigrams fr-alice-partial", args: fmt.Sprintf("-w -a fr -s 3 -o %s %s", outPath, inputFRAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputFRAliceW3)
		}},

		// Discover

		{desc: "discover fr", args: fmt.Sprintf("-d -o %s %s", outPath, inputFRAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)

			langs, err := alphabet.LoadLanguagesFromFile(outPath)
			require.NoError(t, err)
			lang, err := langs.Get("unknown")
			require.NoError(t, err)
			assert.Equal(t, "unknown", lang.Name)

			assert.ElementsMatch(t, []rune(`!"'()*,-.:;?[]_abcdefghijlmnopqrstuvxyzàâçèéêîôùûœ`), []rune(lang.Letters))
		}},

		// Update

		{desc: "update af", args: fmt.Sprintf("-u -a af -s 2 -o %s %s", outPath, inputAFControl), testFunc: func(t *testing.T) {
			os.Remove(outPath)
			require.NoFileExists(t, outPath)

			_, stdErr, err := runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)

			ft, err := ngrams.LoadFrequenciesFromFile(outPath)
			require.NoError(t, err)

			freq, exists := ft.Get("ôr")
			assert.True(t, exists)

			beforeUpdate := freq.Count

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			_, stdErr, err = runMain()
			require.NoError(t, err)
			assert.Empty(t, stdErr)

			ft, err = ngrams.LoadFrequenciesFromFile(outPath)
			require.NoError(t, err)

			freq, exists = ft.Get("ôr")
			assert.True(t, exists)

			assert.Equal(t, beforeUpdate*2, freq.Count)
		}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Fake CLI args for flag package
			os.Args = make([]string, 0)
			os.Args = append(os.Args, "ngrams")
			os.Args = append(os.Args, strings.Split(tC.args, " ")...)

			tC.testFunc(t)

			// Reset flag package for next test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		})
	}
}

func runMain() (string, string, error) {
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	err := app.Main(&outBuffer, &errBuffer)
	return outBuffer.String(), errBuffer.String(), err
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

func tempOutputPath() string {
	// //AJ### DEBUG
	// return filepath.Join("/Users/andre/temp", "ngrams-unit-testing-out.csv")
	return filepath.Join(os.TempDir(), "ngrams-unit-testing-out.csv")
}

func compareTwoFrequencyTableFiles(t *testing.T, a string, b string) {
	aft, err := ngrams.LoadFrequenciesFromFile(a)
	require.NoError(t, err)

	bft, err := ngrams.LoadFrequenciesFromFile(b)
	require.NoError(t, err)

	ae := aft.EntriesSortedByCount()
	be := bft.EntriesSortedByCount()
	assert.Equal(t, len(ae), len(be), "not the same number of rows")

	for i, af := range ae {
		bf := be[i]
		assert.Equal(t, af.Token, bf.Token)
		assert.Equal(t, af.Count, bf.Count)
		assert.InEpsilon(t, af.Percentage, bf.Percentage, 0.000000001)
	}
}

const (
	ngramTestData = "../../../text/ngrams/testdata/"

	inputENControl = ngramTestData + "en-control.txt"
	inputAFControl = ngramTestData + "af-control.txt"
	inputENAlice   = ngramTestData + "en-alice-partial.txt"
	inputFRAlice   = ngramTestData + "fr-alice-partial.txt"

	// Letters
	outputENControl1 = ngramTestData + "freq-1-en-control.csv"
	outputENControl2 = ngramTestData + "freq-2-en-control.csv"

	outputAFControl1 = ngramTestData + "freq-1-af-control.csv"
	outputAFControl2 = ngramTestData + "freq-2-af-control.csv"

	outputENAlice2 = ngramTestData + "freq-2-en-alice.csv"
	outputENAlice3 = ngramTestData + "freq-3-en-alice.csv"

	outputFRAlice2 = ngramTestData + "freq-2-fr-alice.csv"
	outputFRAlice3 = ngramTestData + "freq-3-fr-alice.csv"

	// Words
	outputAFControlW2 = ngramTestData + "freq-2w-af-control.csv"

	outputENAliceW2 = ngramTestData + "freq-2w-en-alice.csv"
	outputENAliceW3 = ngramTestData + "freq-3w-en-alice.csv"

	outputFRAliceW2 = ngramTestData + "freq-2w-fr-alice.csv"
	outputFRAliceW3 = ngramTestData + "freq-3w-fr-alice.csv"
)
