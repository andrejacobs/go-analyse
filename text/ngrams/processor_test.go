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

package ngrams_test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-analyse/text/ngrams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessorLoadFrequenciesFromFile(t *testing.T) {
	p := ngrams.NewFrequencyProcessor(ngrams.ProcessWords, alphabet.MustBuiltin("en"), 1)
	err := p.LoadFrequenciesFromFile("testdata/freq-load-test.txt")
	require.NoError(t, err)

	ft := p.FrequencyTable()
	assert.Equal(t, 3, ft.Len())
	assert.ElementsMatch(t, []string{"the", "fox", "she"}, ft.Tokens())
}

func TestProcessorLoadFrequenciesFromFileFail(t *testing.T) {
	p := ngrams.NewFrequencyProcessor(ngrams.ProcessWords, alphabet.MustBuiltin("en"), 1)
	err := p.LoadFrequenciesFromFile("testdata/freq-load-fail.txt")
	assert.ErrorContains(t, err, "failed to parse the count field from the csv")
}

func TestProcessorLoadAndSaveFrequenciesFromFile(t *testing.T) {
	p := ngrams.NewFrequencyProcessor(ngrams.ProcessWords, alphabet.MustBuiltin("en"), 1)
	err := p.LoadFrequenciesFromFile("testdata/freq-load-test.txt")
	require.NoError(t, err)

	temp := filepath.Join(os.TempDir(), "ngrams-unit-test.csv")
	defer os.Remove(temp)
	require.NoError(t, p.Save(temp))

	p2 := ngrams.NewFrequencyProcessor(ngrams.ProcessWords, alphabet.MustBuiltin("en"), 1)
	require.NoError(t, p2.LoadFrequenciesFromFile(temp))

	compareTwoFrequencyTables(t, p.FrequencyTable(), p2.FrequencyTable())
}

func TestProcessorProcessFiles(t *testing.T) {
	testCases := []struct {
		desc      string
		paths     []string
		lang      alphabet.Language
		tokenSize int
		words     bool
		errMsg    string
		expFreqs  string
		testFunc  func(t *testing.T, p *ngrams.FrequencyProcessor)
	}{
		// Letters
		{desc: "af-control 1",
			paths: []string{"testdata/af-control.txt"},
			lang:  alphabet.MustBuiltin("af"), tokenSize: 1,
			expFreqs: "testdata/freq-1-af-control.csv",
		},
		{desc: "af-control 2",
			paths: []string{"testdata/af-control.txt"},
			lang:  alphabet.MustBuiltin("af"), tokenSize: 2,
			expFreqs: "testdata/freq-2-af-control.csv",
		},
		{desc: "en-control 1",
			paths: []string{"testdata/en-control.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 1,
			expFreqs: "testdata/freq-1-en-control.csv",
		},
		{desc: "en-control 2",
			paths: []string{"testdata/en-control.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 2,
			expFreqs: "testdata/freq-2-en-control.csv",
		},

		{desc: "en-alice-partial 1",
			paths: []string{"testdata/en-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 1,
			expFreqs: "testdata/freq-1-en-alice.csv",
		},
		{desc: "en-alice-partial 2",
			paths: []string{"testdata/en-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 2,
			expFreqs: "testdata/freq-2-en-alice.csv",
		},
		{desc: "en-alice-partial 3",
			paths: []string{"testdata/en-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 3,
			expFreqs: "testdata/freq-3-en-alice.csv",
		},

		{desc: "fr-alice-partial 1",
			paths: []string{"testdata/fr-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("fr"), tokenSize: 1,
			expFreqs: "testdata/freq-1-fr-alice.csv",
		},
		{desc: "fr-alice-partial 2",
			paths: []string{"testdata/fr-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("fr"), tokenSize: 2,
			expFreqs: "testdata/freq-2-fr-alice.csv",
		},
		{desc: "fr-alice-partial 3",
			paths: []string{"testdata/fr-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("fr"), tokenSize: 3,
			expFreqs: "testdata/freq-3-fr-alice.csv",
		},

		// Words
		{desc: "af-control words 1",
			paths: []string{"testdata/af-control.txt"},
			lang:  alphabet.MustBuiltin("af"), tokenSize: 1, words: true,
			expFreqs: "testdata/freq-1w-af-control.csv",
		},
		{desc: "af-control words 2",
			paths: []string{"testdata/af-control.txt"},
			lang:  alphabet.MustBuiltin("af"), tokenSize: 2, words: true,
			expFreqs: "testdata/freq-2w-af-control.csv",
		},
		{desc: "en-control words 1",
			paths: []string{"testdata/en-control.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 1, words: true,
			expFreqs: "testdata/freq-1w-en-control.csv",
		},
		{desc: "en-control words 2",
			paths: []string{"testdata/en-control.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 2, words: true,
			expFreqs: "testdata/freq-2w-en-control.csv",
		},

		{desc: "en-alice-partial words 1",
			paths: []string{"testdata/en-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 1, words: true,
			expFreqs: "testdata/freq-1w-en-alice.csv",
		},
		{desc: "en-alice-partial words 2",
			paths: []string{"testdata/en-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 2, words: true,
			expFreqs: "testdata/freq-2w-en-alice.csv",
		},
		{desc: "en-alice-partial words 3",
			paths: []string{"testdata/en-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 3, words: true,
			expFreqs: "testdata/freq-3w-en-alice.csv",
		},

		{desc: "fr-alice-partial words 1",
			paths: []string{"testdata/fr-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("fr"), tokenSize: 1, words: true,
			expFreqs: "testdata/freq-1w-fr-alice.csv",
		},
		{desc: "fr-alice-partial words 2",
			paths: []string{"testdata/fr-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("fr"), tokenSize: 2, words: true,
			expFreqs: "testdata/freq-2w-fr-alice.csv",
		},
		{desc: "fr-alice-partial words 3",
			paths: []string{"testdata/fr-alice-partial.txt"},
			lang:  alphabet.MustBuiltin("fr"), tokenSize: 3, words: true,
			expFreqs: "testdata/freq-3w-fr-alice.csv",
		},

		{desc: "multiple files words 1",
			paths: []string{"testdata/af-control.txt", "testdata/en-control.txt"},
			lang:  alphabet.MustBuiltin("af"), tokenSize: 1, words: true,
			testFunc: func(t *testing.T, p *ngrams.FrequencyProcessor) {
				ft, err := loadWordFrequenciesFromFiles([]string{
					"testdata/en-control.txt",
					"testdata/af-control.txt"},
					alphabet.MustBuiltin("af"), 1)
				require.NoError(t, err)
				compareTwoFrequencyTables(t, ft, p.FrequencyTable())
			},
		},

		// Zip file support
		{desc: "zip file",
			paths: []string{"testdata/collection1.zip"},
			lang:  alphabet.MustBuiltin("en"), tokenSize: 1, words: true,
			testFunc: func(t *testing.T, p *ngrams.FrequencyProcessor) {
				ft, err := loadWordFrequenciesFromFiles([]string{
					"testdata/en-control.txt",
					"testdata/af-control.txt",
					"testdata/en-alice-partial.txt",
					"testdata/fr-alice-partial.txt"},
					alphabet.MustBuiltin("en"), 1)
				require.NoError(t, err)
				compareTwoFrequencyTables(t, ft, p.FrequencyTable())
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := ngrams.NewFrequencyProcessor(ngrams.ProcessorMode(tC.words), tC.lang, tC.tokenSize)
			err := p.ProcessFiles(context.Background(), tC.paths)

			if tC.errMsg != "" {
				assert.ErrorContains(t, err, tC.errMsg)
			}

			if tC.expFreqs != "" {
				expected, err := ngrams.LoadFrequenciesFromFile(tC.expFreqs)
				require.NoError(t, err)
				compareTwoFrequencyTables(t, expected, p.FrequencyTable())
			}

			if tC.testFunc != nil {
				tC.testFunc(t, p)
			}
		})
	}
}

//-----------------------------------------------------------------------------

func loadWordFrequenciesFromFiles(paths []string, language alphabet.Language, tokenSize int) (*ngrams.FrequencyTable, error) {
	ft := ngrams.NewFrequencyTable()

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open the file %q. %w", path, err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: failed to close %s. %v", path, err)
			}
		}()

		r := bufio.NewReader(f)

		if err := ft.ParseWordTokens(context.Background(), r, language, tokenSize); err != nil {
			return nil, err
		}
	}

	ft.Update()
	return ft, nil
}
