package alphabet_test

import (
	"fmt"
	"testing"

	"github.com/andrejacobs/go-analysis/internal/alphabet"
	"github.com/stretchr/testify/assert"
)

func TestGeneratedLanguages(t *testing.T) {
	// Simple unit-test just to ensure the generator was run and produces
	// some of the expected output
	assert.Equal(t, alphabet.Languages()["af"], alphabet.Language{Name: "Afrikaans", Code: "af", Letters: "abcdefghijklmnopqrstuvwxyzáêéèëïíîôóúû"})
	assert.Equal(t, alphabet.Languages()["en"], alphabet.Language{Name: "English", Code: "en", Letters: "abcdefghijklmnopqrstuvwxyz"})
	assert.Equal(t, alphabet.Languages()["es"], alphabet.Language{Name: "Spanish", Code: "es", Letters: "abcdefghijklmnopqrstuvwxyzáéíñóúü"})
	assert.Equal(t, alphabet.Languages()["da"], alphabet.Language{Name: "Danish", Code: "da", Letters: "abcdefghijklmnopqrstuvwxyzæøå"})
	assert.Equal(t, alphabet.Languages()["ar"], alphabet.Language{Name: "Arabic", Code: "ar", Letters: "أابتثجحخدذرزسشصضطظعغفقكلمنهؤوئىيء"})

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
			assert.Equal(t, tC.expected, alphabet.Languages()[tC.code].ContainsRune(tC.check))
		})
	}
}
