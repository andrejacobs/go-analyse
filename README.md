# go-analysis

Analysis related code written in Go

## Supported languages

Package `internal/alphabet` is used to describe the valid lowercase letters (runes) for a language that we can then
run various analysis on a given corpora.

The default set of languages are generated from the `internal/alphabet/testdata/languages.csv` file and by running
the command `make go-generate`.

## TODO

Disclaimer about not being a linguistic expert. Alphabet in this code means a valid set of Unicode runes that
will be used to determine n-grams. Letter in this code means a Unicode rune (generally no numbers or symbols).
Language means a set of identifyable writing Letters (e.g. EN = English = Latin Alphabet).
