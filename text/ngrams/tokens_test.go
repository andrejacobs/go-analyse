package ngrams_test

import (
	"bufio"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-analyse/text/ngrams"
	"github.com/andrejacobs/go-collection/collection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLetterTokens(t *testing.T) {
	testCases := []struct {
		desc      string
		filename  string
		tokenSize int
		language  alphabet.Language
		expected  collection.Set[string]
		//TODO: support context cancel, with check
		// raise an error
		// set of ngrams check
	}{
		{
			desc:      "Monogram - Alice",
			filename:  "testdata/en-alice-partial.txt",
			tokenSize: 1,
			language:  alphabet.Languages()["en"],
			expected:  collection.NewSetFrom(strings.Split(alphabet.Languages()["en"].Letters, "")),
		},
		{
			desc:      "Bigram - Alice",
			filename:  "testdata/en-alice-partial.txt",
			tokenSize: 2,
			language:  alphabet.Languages()["en"],
			expected:  collection.NewSetFrom([]string{"ae"}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			f, err := os.Open(tC.filename)
			require.NoError(t, err)
			defer f.Close()

			ctx := context.Background()
			r := bufio.NewReader(f)

			result := collection.NewSet[string]()

			err = ngrams.ParseLetterTokens(ctx, r, tC.language, tC.tokenSize,
				func(token string, err error) error {
					assert.Len(t, token, tC.tokenSize)
					result.Insert(token)
					return nil
				})
			require.NoError(t, err)

			assert.Equal(t, tC.expected, result)
		})
	}
}

//TODO: Add an expected file that we load the set from
// bigrams-en-alice-partial.txt
