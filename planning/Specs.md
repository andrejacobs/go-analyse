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

-   `cmd1`

    -   Does A
    -   Does B
    -   etc.

-   `-h, --help`

    -   Displays usage and help information.

-   `-v, --version`
    -   Displays the version of the tool.

### Bonus features

-   `cmd2`

    -   Does A

## Module summary

Overview of the packages provided by this module

-   `text/ngrams`

    -   Calculate monograms (1-gram), bigrams (2-grams), trigrams (3-grams), quadgrams (4-grams) and
        quintgrams (5-grams) for either letters or words on a given corpora.

-   `text/alphabet`
    -   Describe the valid set of unicode runes that are used in other packages to perform text analysis.

## Reference

-   Lib1
-   Tool1
-   etc.
