package app_test

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/cmd/ngrams/app"
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
