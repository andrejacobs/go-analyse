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

	_, err = alphabet.Builtin("golang")
	assert.ErrorContains(t, err, "no built-in language found with code \"golang\"")

	assert.NotEmpty(t, alphabet.BuiltinLanguages())

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
