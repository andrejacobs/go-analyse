package app_test

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/cmd/ngrams/app"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	testCases := []struct {
		desc     string
		args     string
		testFunc func(t *testing.T)
	}{
		{desc: "no input args", args: "", testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, "failed to configure the app. expected at least one input path")
			assert.Error(t, err)
		}},

		{desc: "invalid size", args: "-s 0 ./in.txt", testFunc: func(t *testing.T) {
			_, stdErr, err := runMain()
			assert.Contains(t, stdErr, "invalid ngram size 0")
			assert.Error(t, err)
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
