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

package alphabet

import (
	"fmt"
	"strings"
)

//go:generate go run generate_languages.go

// LanguageCode describes an ISO 639 set 1 language code.
type LanguageCode string

// Language describes the alphabet letters found in a language.
type Language struct {
	// Name of the language (e.g. Afrikaans).
	Name string
	// ISO 639 set 1 language code (e.g. af) https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes.
	Code LanguageCode
	// Letters (in UTF-8 and in lowercase) found in the language.
	Letters string
}

// LanguageMap is used to map from a language code to info about the language.
type LanguageMap map[LanguageCode]Language

// ContainsRune returns true if the language contains the rune.
// The letters of the language is expected to only contain the lowercase runes
// that make up the alphabet and thus the specified rune is assumed to be a lowercase rune as well.
func (l Language) ContainsRune(r rune) bool {
	return strings.ContainsRune(l.Letters, r)
}

// Get the language for the given code or return an error.
func (lm LanguageMap) Get(code LanguageCode) (Language, error) {
	lang, exists := lm[code]
	if !exists {
		return Language{}, fmt.Errorf("no language found with code %q", code)
	}
	return lang, nil
}
