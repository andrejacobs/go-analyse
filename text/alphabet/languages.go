// DO NOT EDIT. This code is generated by generate_languages.go

package alphabet

import (
	"fmt"
)

var languages = LanguageMap{
	"af": Language{Name: "Afrikaans", Code: "af", Letters: "abcdefghijklmnopqrstuvwxyzáêéèëïíîôóúû"},
	"ar": Language{Name: "Arabic", Code: "ar", Letters: "أابتثجحخدذرزسشصضطظعغفقكلمنهؤوئىيء"},
	"da": Language{Name: "Danish", Code: "da", Letters: "abcdefghijklmnopqrstuvwxyzæøå"},
	"de": Language{Name: "German", Code: "de", Letters: "abcdefghijklmnopqrstuvwxyzäöüß"},
	"en": Language{Name: "English", Code: "en", Letters: "abcdefghijklmnopqrstuvwxyz"},
	"es": Language{Name: "Spanish", Code: "es", Letters: "abcdefghijklmnopqrstuvwxyzáéíñóúü"},
	"et": Language{Name: "Estonian", Code: "et", Letters: "abcdefghijklmnopqrstuvwxyzäöõü"},
	"fi": Language{Name: "Finnish", Code: "fi", Letters: "abcdefghijklmnopqrstuvwxyzäö"},
	"fr": Language{Name: "French", Code: "fr", Letters: "abcdefghijklmnopqrstuvwxyzàâæçéèêëîïôœùûüÿ"},
	"nl": Language{Name: "Dutch", Code: "nl", Letters: "abcdefghijklmnopqrstuvwxyzàäèéëïĳöü"},
	"sv": Language{Name: "Swedish", Code: "sv", Letters: "abcdefghijklmnopqrstuvwxyzåäö"},
}

// Builtin returns the built-in language for the given ISO 639 set 1 language.
func Builtin(code LanguageCode) (Language, error) {
	lang, exists := languages[code]
	if !exists {
		return Language{}, fmt.Errorf("no built-in language found with code %q", code)
	}
	return lang, nil
}

// MustBuiltin returns the built-in language for the given ISO 639 set 1 language or panics.
func MustBuiltin(code LanguageCode) Language {
	lang, exists := languages[code]
	if !exists {
		panic(fmt.Errorf("no built-in language found with code %q", code))
	}
	return lang
}

// BuiltinLanguages return the built-in languages.
func BuiltinLanguages() LanguageMap {
	return languages
}
