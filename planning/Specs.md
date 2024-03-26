# Specification

What is this module for?

go-analyse is a collection of go packages I use to perform various analysis for other projects of mine. The primary
focus at the moment is just on text analysis (e.g. n-grams).

Who is it for?

The primary audience is just myself for the moment and with the hope that this might be useful to others as well.

Why build this yourself?

Easy, I build this myself because I want to learn more about certain topics. Aka procastination by learning and coding :-D

What prompted the creation of this module?

Sheepishly I will admit that my original goal was to just try and do daily practise of n-gram typing on a new split
keyboard, but why use an online tool when you can write you own in Go and for the terminal. Also why stop there when
you can write your own n-gram generation tools instead of just downloading a data set?
For reasons 1) to procastinate, 2) to learn, 3) to code and 4) to get the dopamine fix.

## Command summary

Overview of the commands that will be available.

### Core features

-   `ngrams`

    -   Create either letter or word ngram frequency tables from the input corpora.
    -   Can update existing frequency tables.
    -   Can discover the alphabet letters used from the input corpora.

    -   `-a, --lang` The alphabet to use. This can either be a built-in (e.g. en = english) or a .csv file.
    -   `--languages` CSV file containing the languages.
    -   `-l, --letters` Create ngrams using letter combinations. E.g. bigrams like st, er, ae
    -   `-w, --words` Create ngrams using word combinations. E.g. bigrams like "the cat", "he jumped"
    -   `-s, --size` ngram size. E.g. 1 = monogram, 2 = bigrams, 3 = trigrams etc.
    -   `-o, --out` The file path to write the results to. (csv format). See `--discover` for when this file is a language file.
    -   `-d, --discover` Instead of creating the ngram frequency table this will discover the non whitespace characters used and write a language file to `--out` path.
    -   `-u, --update` Load the frequency table specified by `--out` and then update as new data is parsed by the corpora.
    -   `--available` Lists the available languages. If --languages file is specified then list out those languages else the built-in ones.

    -   `-h, --help`

        -   Displays usage and help information.

    -   `--version`
        -   Displays the version of the tool.

What options work with what?

-   `-a, --lang`, `--languages`, `-l, -w`, `-s` works with creating the ngram frequency table. Depends on the `-o` and inputs paths.
-   `-u, --update` works with the first set, it first reads the `-o` file and then proceeds to update the frequency table before writing back to `-o`
-   `-d, --discover` works only with the input paths and writes the language file to `-o`.

### Bonus features

-   `ngrams`

    -   Input files can also be zip.
    -   URLs can be specified instead of files to fetch corpora from the web. E.g. A GET request is made to the URL and then parsed.
        -   Would then need to think about allowing netscape style cookies.txt and also setting the user-agent.
        -   What about retries? exponential back-off.
    -   Progress bar `--progress`
    -   Verbose `-v, --verbose`

## Module summary

Overview of the packages provided by this module

-   `text/ngrams`

    -   Calculate monograms (1-gram), bigrams (2-grams), trigrams (3-grams), quadgrams (4-grams) and
        quintgrams (5-grams) for either letters or words on a given corpora.

-   `text/alphabet`
    -   Describe the valid set of unicode runes that are used in other packages to perform text analysis.
