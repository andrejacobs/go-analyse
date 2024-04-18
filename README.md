# go-analyse

Analysis related code written in Go

Please see https://github.com/andrejacobs/datasets-text for example corpora and pre-generated output produced by the tools from this repository.

## Install from source

1. Ensure you have at least Go 1.22 installed: https://go.dev/dl/
1. Clone the repository

    ```
    git clone git@github.com:andrejacobs/go-analyse.git
    ```

1. Build

    ```
    cd go-analyse/
    make build

    # Executables are created at: build/bin/<OS>/<CPU_ARCH> and
    # symlinked as: build/bin/<CLI>
    ```

1. Install

    ```
    make install

    # Executables are copied to: $(go env GOPATH)/bin
    # Ensure $(go env GOPATH)/bin is in your $PATH
    ```

## N-grams

The `ngrams` CLI app can be used to generate letter and word ngram frequency tables from a set of input corpora.

General usage:

```
ngrams [options] [-o output] file ...
```

Zip files are also supported as input files.

See `ngrams --help` for more details on the supported options.

### Examples:

Basic usage:

```
# Calculate the letter bigram frequency for Afrikaans (af)

$ ngrams --size 2 --lang af af-corpus.zip
# produces the output file: af-letters-2.csv



# Calculate the word bigram frequency for Afrikaans (af)

$ ngrams --words --size 2 --lang af af-corpus.zip
# produces the output file: af-words-2.csv



# Specify the output file

$ ngrams --size 2 --lang af --out /path/to/output.csv af-corpus.zip
```

With a progress bar:

```
$ ngrams --progress --words --size 2 --lang af af-corpus.zip
[1/1]  47% |████████        | (49/103 MB, 19 MB/s) [2s:2s]
```

With more verbose information:

```
$ ngrams --verbose --words --size 2 --lang af af-corpus.zip
Language: af - Afrikaans
Generating 2 word ngrams...
[1/1] af-corpus.zip
Saving frequency table...
Created frequency table at: "./af-words-2.csv"
```

To discover the non-whitespace lowercased letters (unicode runes) used in a corpus:

```
$ ngrams --discover corpus1.zip corpus2.zip samples.txt
# produces the output file: languages.csv

$ cat languages.csv
#code,name,letters
unknown,unknown,"!""$%&'()*+,-./0123456789:;<>?@`abcdefghijklmnopqrstuvwxyz~£¥¦§¨«¬®°±²³´µ¶·¹º»¼½¾¿×àáâãäåçèéêëìíîïñòóôõö÷øùúûüýÿďėęīķńŉōőšűȇȏʼ˚˜́̈̓́ίαβγδεικπρςσυωόавгеийклнрсчюёїў٪٬აბდევიკლნრსუქყệὶ‐‑‒–—―‘’‚“”„‡•…‰′″‹›⁰€ℓ™−╔╗♦地憂抧抯揕揟揤梐梕梘梞梟梥梬棐璴璶璼痴醤鉶鑞閚阹雃雐雔雗雛雝飀飊飏飔馻骹黮"%
```

To see the list of available languages:

```
# Built-in languages

$ ngrams --available
af : Afrikaans
ar : Arabic
da : Danish
de : German
en : English
es : Spanish
et : Estonian
fi : Finnish
fr : French
nl : Dutch
sv : Swedish

# Supply a custom languages files

$ ngrams --available --languages languages.csv
af : Afrikaans
golang: Go Programming Language
```

## Supported languages

Package `text/alphabet` is used to describe the valid lowercase letters (runes) for a language that can then
be used to run various analysis on for a given corpora.

The default set of languages are generated from the `text/alphabet/testdata/languages.csv` file and by running
the command `make go-generate`.

## Packages

### `text/alphabet`

Provides support to describe a set of languages and the alphabets used in the writing style for that language.

Built-in languages: (provided by "text/alphabet/languages.go")

```go
lang, err := alphabet.Builtin("af")
print(lang.Name) // Afrikaans

lang := alphabet.MustBuiltin("af") // Will panic if the language does not exist

// The map of language code (generally the ISO 639 set 1 code) to the Language struct
languages := alphabet.BuiltinLanguages()
```

To update the built-in languages:

1. Update the "testdata/languages.csv" file.
2. Run `make go-generate` which will then create the file "text/alphabet/languages.go"

Load a languages file:

example.csv

```
#code,name,letters
af,Afrikaans,abcdefghijklmnopqrstuvwxyzáêéèëïíîôóúû
en,English,abcdefghijklmnopqrstuvwxyz
```

```go
languages, err := alphabet.LoadLanguagesFromFile("example.csv")
```

To discover the unique unicode runes used in files:

```go
p := alphabet.NewDiscoverProcessor()
err := p.ProcessFiles(context.Background(), []string{"discover1.txt", "example2.txt"})

// Save the discovered runes in the supported CSV format
err = p.Save("example.csv")

languages, err := alphabet.LoadLanguagesFromFile("example.csv")
```

### `text/ngrams`

Provides support to calculate the frequency tables for letter and word n-grams.
Examples of ngrams:

-   letters

    -   monograms: a, b, c, d
    -   bigrams: he, sh, th, ng, in
    -   trigrams: the, fox

-   words
    -   monograms: apple, the, quick
    -   bigrams: the quick, fox jumped
    -   trigrams: the quick brown, fox jumped over

To calculate the n-grams from input files:

```go
// Word bigrams using the built-in "en" language
p := ngrams.NewFrequencyProcessor(ngrams.ProcessWords, alphabet.MustBuiltin("en"), 2)
err := p.LoadFrequenciesFromFile("input-corpus.txt")

// The calculated frequency table
ft := p.FrequencyTable()

// To save the frequency table
err = p.Save("en-word-bigrams.csv")
```

Frequency table file format in CSV

```
#token,count,percentage
the,5,0.1
fox,2,0.03
```

## Glossary

This section describes in general the words used and the meaning in the context of this code repository.

-   alphabet: a valid set of unicode runes that describes the writing letters used in a language.
-   language: a set of identifyable writing letters that describes a language. Identified by an ISO 639 set 1 code (e.g. EN = English).
-   letter: a lowercased unicode rune (generally no numbers or symbols).
