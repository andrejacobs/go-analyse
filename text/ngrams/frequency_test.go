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
	defer os.Remove(f.Name())

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

func TestFrequencyEntriesSortedByCount(t *testing.T) {
	freq := ngrams.NewFrequencyTable()
	freq.Add("a", 1)
	freq.Add("b", 2)
	freq.Add("c", 1)
	freq.Add("d", 1)
	freq.Add("e", 1)
	freq.Add("f", 2)

	expected := []ngrams.Frequency{
		{Token: "b", Count: 2},
		{Token: "f", Count: 2},
		{Token: "a", Count: 1},
		{Token: "c", Count: 1},
		{Token: "d", Count: 1},
		{Token: "e", Count: 1},
	}

	// Should always produce the exact same sort order
	for i := 0; i < 100; i++ {
		assert.Equal(t, expected, freq.EntriesSortedByCount())
	}
}

func compareTwoFrequencyTables(t *testing.T, a *ngrams.FrequencyTable, b *ngrams.FrequencyTable) {
	ae := a.EntriesSortedByCount()
	be := b.EntriesSortedByCount()
	assert.Equal(t, len(ae), len(be), "not the same number of rows")

	for i, af := range ae {
		bf := be[i]
		assert.Equal(t, af.Token, bf.Token)
		assert.Equal(t, af.Count, bf.Count)
		assert.InEpsilon(t, af.Percentage, bf.Percentage, 0.001,
			"token a: %s, %d, %g\ntoken b: %s, %d, %g",
			af.Token, af.Count, af.Percentage, bf.Token, bf.Count, bf.Percentage)
	}
}
