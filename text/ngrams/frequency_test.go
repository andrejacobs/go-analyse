package ngrams_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analysis/text/ngrams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFrequencies(t *testing.T) {
	r := strings.NewReader(`#token,count,percentage
the,2,0.3
fox,1,0.01
#ignore,,


she,10,0.9    
#overwrite,,
the,5,0.7
`)

	freq, err := ngrams.LoadFrequencies(r)
	require.NoError(t, err)

	assert.Equal(t, 3, freq.Len())
	assert.ElementsMatch(t, []string{"the", "fox", "she"}, freq.Tokens())

	a, ok := freq.Get("she")
	assert.True(t, ok)
	b, ok := freq.Get("the")
	assert.True(t, ok)
	c, ok := freq.Get("fox")
	assert.True(t, ok)

	assert.Equal(t, []ngrams.Frequency{a, b, c}, freq.EntriesSortedByCount())
}

func TestLoadFrequenciesFromFile(t *testing.T) {
	freq, err := ngrams.LoadFrequenciesFromFile("testdata/freq-load-test.txt")
	require.NoError(t, err)

	assert.Equal(t, 3, freq.Len())
	assert.ElementsMatch(t, []string{"the", "fox", "she"}, freq.Tokens())

	a, ok := freq.Get("she")
	assert.True(t, ok)
	b, ok := freq.Get("the")
	assert.True(t, ok)
	c, ok := freq.Get("fox")
	assert.True(t, ok)

	assert.Equal(t, []ngrams.Frequency{a, b, c}, freq.EntriesSortedByCount())
}

func TestLoadFrequenciesParseErrors(t *testing.T) {
	testCases := []struct {
		input  string
		errMsg string
	}{
		{input: "the,nan,0.1", errMsg: "failed to parse the count field from the csv"},
		{input: "the,42,abc", errMsg: "failed to parse the percentage field from the csv"},
	}
	for i, tC := range testCases {
		t.Run(fmt.Sprintf("ParseError%d", i), func(t *testing.T) {
			r := strings.NewReader(tC.input)
			_, err := ngrams.LoadFrequencies(r)
			assert.ErrorContains(t, err, tC.errMsg)
		})
	}
}
