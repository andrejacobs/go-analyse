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
		{desc: "af-control 1",
			paths: []string{"testdata/af-control.txt"},
			lang:  alphabet.MustBuiltin("af"), tokenSize: 1,
			expFreqs: "testdata/freq-1-af-control.csv",
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
