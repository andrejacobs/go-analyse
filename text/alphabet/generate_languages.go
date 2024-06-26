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

//go:build ignore
// +build ignore

// Generates the languages.go file

package main

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"golang.org/x/exp/maps"
)

const (
	outputFilename = "languages.go"
	inputData      = "testdata/languages.csv"
)

func main() {
	f, err := os.Create(outputFilename)
	if err != nil {
		die(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			die(err)
		}
	}()

	fmt.Printf("Generating %s\n", outputFilename)

	if err := writeHeader(f); err != nil {
		die(err)
	}

	if err := processCSV(f, inputData); err != nil {
		die(err)
	}

	if err := writeFooter(f); err != nil {
		die(err)
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func writeHeader(w io.Writer) error {
	const header = `// DO NOT EDIT. This code is generated by generate_languages.go

package alphabet

import (
	"fmt"
)

var languages = LanguageMap{
`

	_, err := io.WriteString(w, header)
	if err != nil {
		return fmt.Errorf("failed to write the header. %w", err)
	}

	return nil
}

func writeFooter(w io.Writer) error {
	const footer = `
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

`
	_, err := io.WriteString(w, footer)
	if err != nil {
		return fmt.Errorf("failed to write the footer. %w", err)
	}

	return nil
}

func processCSV(w io.Writer, path string) error {
	languages, err := alphabet.LoadLanguagesFromFile(path)
	if err != nil {
		return err
	}

	keys := maps.Keys(languages)
	slices.Sort(keys)

	for _, code := range keys {
		lang := languages[code]

		io.WriteString(w, "\t"+fmt.Sprintf(`"%s": Language{Name: "%s", Code: "%s", Letters: "%s"},`+"\n",
			code, lang.Name, lang.Code, lang.Letters))
	}

	return nil
}
