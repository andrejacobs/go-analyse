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
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	_ "golang.org/x/exp/maps"
)

// ErrNoLanguages is returned when no languages could be loaded.
var ErrNoLanguages = errors.New("no languages")

// LoadLanguages parses a set of languages from an io.Reader.
//
// Expected CSV format in UTF-8: code,name,letters
// Lines starting with a # is ignored.
func LoadLanguages(r io.Reader) (LanguageMap, error) {
	result := make(LanguageMap)
	csvR := csv.NewReader(r)

	for {
		record, err := csvR.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse csv. %w", err)
		}

		if len(record) < 3 {
			continue
		}

		if strings.HasPrefix(record[0], "#") {
			continue
		}

		code := LanguageCode(record[0])
		l := Language{
			Name:    record[1],
			Code:    code,
			Letters: strings.ToLower(record[2]),
		}

		result[code] = l
	}

	if len(result) < 1 {
		return nil, ErrNoLanguages
	}

	return result, nil
}

// LoadLanguagesFromFile parses a set of languages from a UTF-8 encoded text file.
// See [LoadLanguages] for more details.
func LoadLanguagesFromFile(path string) (LanguageMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close %q. %v", path, err)
		}
	}()

	result, err := LoadLanguages(f)
	if err != nil {
		return nil, fmt.Errorf("failed to load languages from %q. %w", path, err)
	}

	return result, nil
}
