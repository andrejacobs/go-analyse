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
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-analyse/text/ngrams"
	"github.com/andrejacobs/go-collection/collection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//AJ### TODO: Review this. Might be redundant testing given the helpers_test

func TestParseLetterTokens(t *testing.T) {
	enLang, err := alphabet.Builtin("en")
	require.NoError(t, err)
	afLang, err := alphabet.Builtin("af")
	require.NoError(t, err)

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
			desc:      "EN - Control Monograms",
			filename:  "testdata/en-control.txt",
			tokenSize: 1,
			language:  enLang,
			expected:  "testdata/freq-1-en-control.csv",
		},
		{
			desc:      "EN - Control Bigrams",
			filename:  "testdata/en-control.txt",
			tokenSize: 2,
			language:  enLang,
			expected:  "testdata/freq-2-en-control.csv",
		},
		{
			desc:      "AF - Control Monograms",
			filename:  "testdata/af-control.txt",
			tokenSize: 1,
			language:  afLang,
			expected:  "testdata/freq-1-af-control.csv",
		},
		{
			desc:      "AF - Control Bigrams",
			filename:  "testdata/af-control.txt",
			tokenSize: 2,
			language:  afLang,
			expected:  "testdata/freq-2-af-control.csv",
		},
		{
			desc:      "Monogram - Alice",
			filename:  "testdata/en-alice-partial.txt",
			tokenSize: 1,
			language:  enLang,
			expected:  "testdata/freq-1-en-alice.csv",
		},
		{
			desc:      "Bigram - Alice",
			filename:  "testdata/en-alice-partial.txt",
			tokenSize: 2,
			language:  enLang,
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
					assert.Equal(t, utf8.RuneCountInString(token), tC.tokenSize)
					result.Insert(token)
					return nil
				})
			require.NoError(t, err)

			assert.Equal(t, exp, result)
		})
	}
}

func TestParseWordTokens(t *testing.T) {
	enLang, err := alphabet.Builtin("en")
	require.NoError(t, err)
	afLang, err := alphabet.Builtin("af")
	require.NoError(t, err)

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
			desc:      "EN - Control Monograms",
			filename:  "testdata/en-control.txt",
			tokenSize: 1,
			language:  enLang,
			expected:  "testdata/freq-1w-en-control.csv",
		},
		{
			desc:      "EN - Control Bigrams",
			filename:  "testdata/en-control.txt",
			tokenSize: 2,
			language:  enLang,
			expected:  "testdata/freq-2w-en-control.csv",
		},
		{
			desc:      "AF - Control Monograms",
			filename:  "testdata/af-control.txt",
			tokenSize: 1,
			language:  afLang,
			expected:  "testdata/freq-1w-af-control.csv",
		},
		{
			desc:      "AF - Control Bigrams",
			filename:  "testdata/af-control.txt",
			tokenSize: 2,
			language:  afLang,
			expected:  "testdata/freq-2w-af-control.csv",
		},
		{
			desc:      "Monogram - Alice",
			filename:  "testdata/en-alice-partial.txt",
			tokenSize: 1,
			language:  enLang,
			expected:  "testdata/freq-1w-en-alice.csv",
		},
		{
			desc:      "Bigram - Alice",
			filename:  "testdata/en-alice-partial.txt",
			tokenSize: 2,
			language:  enLang,
			expected:  "testdata/freq-2w-en-alice.csv",
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

			err = ngrams.ParseWordTokens(ctx, r, tC.language, tC.tokenSize,
				func(token string, err error) error {
					assert.Equal(t, tC.tokenSize, len(strings.Split(token, " ")))
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
