package ngrams_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/text/ngrams"
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

func TestFrequenciesLoadAndSave(t *testing.T) {
	freq := ngrams.NewFrequencyTable()
	freq.Add("he", 1)
	freq.Add("he", 2)
	freq.Add("she", 1)
	freq.Add("the", 100)

	f, err := os.CreateTemp("", "unit-testing-ngrams")
	require.NoError(t, err)

	err = freq.Save(f)
	f.Close()
	require.NoError(t, err)

	load, err := ngrams.LoadFrequenciesFromFile(f.Name())
	require.NoError(t, err)
	assert.Equal(t, freq.EntriesSortedByCount(), load.EntriesSortedByCount())
}

func TestFrequencyAdd(t *testing.T) {
	freq := ngrams.NewFrequencyTable()
	freq.Add("he", 1)
	freq.Add("he", 2)
	freq.Add("she", 1)
	freq.Add("the", 100)

	expected := []ngrams.Frequency{
		{Token: "the", Count: 100},
		{Token: "he", Count: 3},
		{Token: "she", Count: 1},
	}

	assert.Equal(t, expected, freq.EntriesSortedByCount())
}

// Use this to generate the test data
// Not ideal to be doing chicken-n-egg and generating the testdata using the code
// you actually want to test. I will just have to be 100% sure of the results
// func TestGenerateFrequencies(t *testing.T) {
// 	genFn := func(input string, output string, lang alphabet.LanguageCode, tokenSize int) {
// 		ft, err := ngrams.FrequencyTableByParsingLetters(context.Background(),
// 			input, alphabet.Languages()[lang], tokenSize)
// 		require.NoError(t, err)
// 		out, err := os.Create(output)
// 		require.NoError(t, err)
// 		ft.Save(out)
// 		out.Close()
// 	}

// 	genFn("testdata/en-alice-partial.txt", "testdata/freq-1-en-alice.csv", "en", 1)
// 	genFn("testdata/en-alice-partial.txt", "testdata/freq-2-en-alice.csv", "en", 2)
// }
