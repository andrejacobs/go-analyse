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

package alphabet_test

import (
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadLanguages(t *testing.T) {
	r := strings.NewReader(`#code,name,letters
en,English,abcdefghijklmnopqrstuvwxyz

weird,Weird,éÅß
coding,Coding,{}[]()/$

# Override previous entry,,
coding,Coding,{}[]/$^%
`)

	languages, err := alphabet.LoadLanguages(r)
	require.NoError(t, err)
	assert.Equal(t, 3, len(languages))

	assert.Equal(t, languages["en"], alphabet.Language{Name: "English", Code: "en", Letters: "abcdefghijklmnopqrstuvwxyz"})
	assert.Equal(t, languages["weird"], alphabet.Language{Name: "Weird", Code: "weird", Letters: "éåß"})
	assert.Equal(t, languages["coding"], alphabet.Language{Name: "Coding", Code: "coding", Letters: `{}[]/$^%`})
}

func TestLoadLanguagesFromFile(t *testing.T) {
	languages, err := alphabet.LoadLanguagesFromFile("testdata/languages.csv")
	require.NoError(t, err)

	assert.Contains(t, languages, alphabet.LanguageCode("en"))
	assert.Contains(t, languages, alphabet.LanguageCode("af"))

	// generated languages.go should be the exact same as loading "testdata/languages.csv"
	for k, v := range languages {
		b, err := alphabet.Builtin(k)
		require.NoError(t, err)

		assert.Equal(t, v, b)
	}

	_, err = alphabet.LoadLanguagesFromFile("testdata/naf.csv")
	require.ErrorContains(t, err, "failed to open \"testdata/naf.csv\"")
}

func TestLoadLanguagesEmpty(t *testing.T) {
	r := strings.NewReader("")

	_, err := alphabet.LoadLanguages(r)
	assert.ErrorIs(t, err, alphabet.ErrNoLanguages)
}

func TestLanguageMapGet(t *testing.T) {
	languages, err := alphabet.LoadLanguagesFromFile("testdata/languages.csv")
	require.NoError(t, err)

	lang, err := languages.Get("af")
	assert.NoError(t, err)
	assert.Equal(t, alphabet.LanguageCode("af"), lang.Code)

	_, err = languages.Get("zu")
	assert.ErrorContains(t, err, "no language found with code \"zu\"")
}

func TestLoadLanguagesReadFailed(t *testing.T) {
	var r FailReader
	_, err := alphabet.LoadLanguages(&r)
	assert.ErrorContains(t, err, "failed to parse csv")
}
