//go:build ignore
// +build ignore

// Generates the languages.go file

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
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

// Languages returns the map of languages
func Languages() LanguageMap {
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
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close %s. %v", path, err)
		}
	}()

	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to process %s. %w", path, err)
		}

		if len(record) < 3 {
			continue
		}

		if strings.HasPrefix(record[0], "#") {
			continue
		}

		io.WriteString(w, "\t"+fmt.Sprintf(`"%s": Language{Name: "%s", Code: "%s", Letters: "%s"},`+"\n",
			record[0], record[1], record[0], strings.ToLower(record[2])))
	}
	return nil
}
