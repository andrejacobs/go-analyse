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

	outPath := tempOutputPath(t)
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

		{desc: "monograms en-control", args: fmt.Sprintf("-s 1 -o %s %s", outPath, inputENControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENControl1)
		}},

		{desc: "bigrams en-control", args: fmt.Sprintf("-s 2 -o %s %s", outPath, inputENControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENControl2)
		}},

		{desc: "monograms af-control", args: fmt.Sprintf("-a af -s 1 -o %s %s", outPath, inputAFControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputAFControl1)
		}},

		{desc: "bigrams af-control", args: fmt.Sprintf("-a af -s 2 -o %s %s", outPath, inputAFControl), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputAFControl2)
		}},

		{desc: "bigrams en-alice-partial", args: fmt.Sprintf("-s 2 -o %s %s", outPath, inputENAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENAlice2)
		}},

		{desc: "trigrams en-alice-partial", args: fmt.Sprintf("-s 3 -o %s %s", outPath, inputENAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputENAlice3)
		}},

		{desc: "bigrams fr-alice-partial", args: fmt.Sprintf("-a fr -s 2 -o %s %s", outPath, inputFRAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputFRAlice2)
		}},

		{desc: "trigrams fr-alice-partial", args: fmt.Sprintf("-a fr -s 3 -o %s %s", outPath, inputFRAlice), testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.NoError(t, err)
			assert.Empty(t, stdErr)
			compareTwoFrequencyTableFiles(t, outPath, outputFRAlice3)
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

func tempOutputPath(t *testing.T) string {
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

	inputENControl   = ngramTestData + "en-control.txt"
	outputENControl1 = ngramTestData + "freq-1-en-control.csv"
	outputENControl2 = ngramTestData + "freq-2-en-control.csv"

	inputAFControl   = ngramTestData + "af-control.txt"
	outputAFControl1 = ngramTestData + "freq-1-af-control.csv"
	outputAFControl2 = ngramTestData + "freq-2-af-control.csv"

	inputENAlice   = ngramTestData + "en-alice-partial.txt"
	outputENAlice2 = ngramTestData + "freq-2-en-alice.csv"
	outputENAlice3 = ngramTestData + "freq-3-en-alice.csv"

	inputFRAlice   = ngramTestData + "fr-alice-partial.txt"
	outputFRAlice2 = ngramTestData + "freq-2-fr-alice.csv"
	outputFRAlice3 = ngramTestData + "freq-3-fr-alice.csv"
)
