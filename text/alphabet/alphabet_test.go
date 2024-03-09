package alphabet_test

import (
	"fmt"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratedLanguages(t *testing.T) {
	// Simple unit-test just to ensure the generator was run and produces
	// some of the expected output
	lang, err := alphabet.Builtin("af")
	require.NoError(t, err)
	assert.Equal(t, lang, alphabet.Language{Name: "Afrikaans", Code: "af", Letters: "abcdefghijklmnopqrstuvwxyzáêéèëïíîôóúû"})

	lang, err = alphabet.Builtin("en")
	require.NoError(t, err)
	assert.Equal(t, lang, alphabet.Language{Name: "English", Code: "en", Letters: "abcdefghijklmnopqrstuvwxyz"})

	lang, err = alphabet.Builtin("es")
	require.NoError(t, err)
	assert.Equal(t, lang, alphabet.Language{Name: "Spanish", Code: "es", Letters: "abcdefghijklmnopqrstuvwxyzáéíñóúü"})

	lang, err = alphabet.Builtin("da")
	require.NoError(t, err)
	assert.Equal(t, lang, alphabet.Language{Name: "Danish", Code: "da", Letters: "abcdefghijklmnopqrstuvwxyzæøå"})

	lang, err = alphabet.Builtin("ar")
	require.NoError(t, err)
	assert.Equal(t, lang, alphabet.Language{Name: "Arabic", Code: "ar", Letters: "أابتثجحخدذرزسشصضطظعغفقكلمنهؤوئىيء"})

	testCases := []struct {
		code     alphabet.LanguageCode
		check    rune
		expected bool
	}{
		{code: "de", check: rune('ß'), expected: true},
		{code: "af", check: rune('ß'), expected: false},
		{code: "da", check: rune('æ'), expected: true},
		{code: "da", check: rune('Æ'), expected: false},
		{code: "ar", check: rune('ض'), expected: true},
	}
	for i, tC := range testCases {
		t.Run(fmt.Sprintf("RuneCheck-%d", i), func(t *testing.T) {
			lang, err := alphabet.Builtin(tC.code)
			require.NoError(t, err)
			assert.Equal(t, tC.expected, lang.ContainsRune(tC.check))
		})
	}
}
