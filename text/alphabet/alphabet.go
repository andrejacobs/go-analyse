package alphabet

import (
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
