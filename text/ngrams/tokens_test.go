package ngrams_test

import (
	"bufio"
	"context"
	"os"
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
		expected  string
		//TODO: support context cancel, with check
		// raise an error
		// set of ngrams check
	}{
		{
			desc:      "Monogram - Alice",
			filename:  "testdata/en-alice-partial.txt",
			tokenSize: 1,
			language:  alphabet.Languages()["en"],
			expected:  "testdata/freq-1-en-alice.csv",
		},
		{
			desc:      "Bigram - Alice",
			filename:  "testdata/en-alice-partial.txt",
			tokenSize: 2,
			language:  alphabet.Languages()["en"],
			expected:  "testdata/freq-2-en-alice.csv",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			exp, err := tokensFromFrequencyFile(tC.expected)
			require.NoError(t, err)

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

			assert.Equal(t, exp, result)
		})
	}
}

func tokensFromFrequencyFile(path string) (collection.Set[string], error) {
	freq, err := ngrams.LoadFrequenciesFromFile(path)
	if err != nil {
		return collection.Set[string]{}, err
	}

	result := collection.NewSetWithCapacity[string](freq.Len())
	for _, e := range freq.Entries() {
		result.Insert(e.Token)
	}

	return result, nil
}
