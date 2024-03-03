// DO NOT EDIT. This code is generated by generate_languages.go

package alphabet

var languages = LanguageMap{
	"af": Language{Name: "Afrikaans", Code: "af", Letters: "abcdefghijklmnopqrstuvwxyzáêéèëïíîôóúû"},
	"ar": Language{Name: "Arabic", Code: "ar", Letters: "أابتثجحخدذرزسشصضطظعغفقكلمنهؤوئىيء"},
	"da": Language{Name: "Danish", Code: "da", Letters: "abcdefghijklmnopqrstuvwxyzæøå"},
	"de": Language{Name: "German", Code: "de", Letters: "abcdefghijklmnopqrstuvwxyzäöüß"},
	"en": Language{Name: "English", Code: "en", Letters: "abcdefghijklmnopqrstuvwxyz"},
	"es": Language{Name: "Spanish", Code: "es", Letters: "abcdefghijklmnopqrstuvwxyzáéíñóúü"},
	"et": Language{Name: "Estonian", Code: "et", Letters: "abcdefghijklmnopqrstuvwxyzäöõü"},
	"fi": Language{Name: "Finnish", Code: "fi", Letters: "abcdefghijklmnopqrstuvwxyzäö"},
	"nl": Language{Name: "Dutch", Code: "nl", Letters: "abcdefghijklmnopqrstuvwxyzàäèéëïĳöü"},
	"sv": Language{Name: "Swedish", Code: "sv", Letters: "abcdefghijklmnopqrstuvwxyzåäö"},
}

// Languages returns the map of languages
func Languages() LanguageMap {
	return languages
}