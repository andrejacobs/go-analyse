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
}
