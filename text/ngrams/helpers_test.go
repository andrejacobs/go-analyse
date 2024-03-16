package ngrams_test

import (
	"context"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-analyse/text/ngrams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateTableByParsingLettersAndWordsFromFiles(t *testing.T) {
	testCases := []struct {
		desc      string
		paths     []string
		lang      alphabet.Language
		tokenSize int
		words     bool
		errMsg    string
		expFreqs  string
		testFunc  func(t *testing.T, ft *ngrams.FrequencyTable)
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

		// Multiple input files
		{desc: "multiple files 1",
			paths: []string{"testdata/af-control.txt", "testdata/en-control.txt"},
			lang:  alphabet.MustBuiltin("af"), tokenSize: 1,
			testFunc: func(t *testing.T, ft *ngrams.FrequencyTable) {
				expected := ngrams.NewFrequencyTable()
				err := expected.UpdateTableByParsingLettersFromFiles(context.Background(),
					[]string{"testdata/en-control.txt", "testdata/af-control.txt"},
					alphabet.MustBuiltin("af"), 1)
				require.NoError(t, err)
				compareTwoFrequencyTables(t, expected, ft)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var err error
			ft := ngrams.NewFrequencyTable()
			if tC.words {
				err = ft.UpdateTableByParsingWordsFromFiles(context.Background(), tC.paths, tC.lang, tC.tokenSize)
			} else {
				err = ft.UpdateTableByParsingLettersFromFiles(context.Background(), tC.paths, tC.lang, tC.tokenSize)
			}

			if tC.errMsg != "" {
				assert.ErrorContains(t, err, tC.errMsg)
			}

			if tC.expFreqs != "" {
				expected, err := ngrams.LoadFrequenciesFromFile(tC.expFreqs)
				require.NoError(t, err)
				compareTwoFrequencyTables(t, expected, ft)
			}

			if tC.testFunc != nil {
				tC.testFunc(t, ft)
			}
		})
	}
}
