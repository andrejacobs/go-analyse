# go-analysis

Analysis related code written in Go

## Supported languages

Package `internal/alphabet` is used to describe the valid lowercase letters (runes) for a language that we can then
run various analysis on a given corpora.

The default set of languages are generated from the `internal/alphabet/testdata/languages.csv` file and by running
the command `make go-generate`.
