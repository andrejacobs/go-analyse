package ngrams_test

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"
	"unicode"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/andrejacobs/go-collection/collection"
	"golang.org/x/exp/maps"
)

func TestWIP(t *testing.T) {
	input := "he SHE 本語本 the sheri zo."

	input = strings.ToLower(input)
	// mono(input, alphabet.Languages()["en"])
	bi(input, alphabet.Languages()["en"])
	tri(input, alphabet.Languages()["en"])
	quad(input, alphabet.Languages()["en"])
}

func mono(input string, lang alphabet.Language) {
	fmt.Printf("monograms: %q\n", input)
	result := make(map[rune]uint64)

	for _, r := range input {
		if lang.ContainsRune(r) {
			count, exists := result[r]
			if !exists {
				result[r] = 1
			} else {
				result[r] = count + 1
			}
		}
	}

	keys := maps.Keys(result)
	slices.Sort(keys)

	for _, k := range keys {
		fmt.Printf("%c = %d\n", k, result[k])
	}

	fmt.Println()
}

func bi(input string, lang alphabet.Language) {
	fmt.Printf("bigrams: %q\n", input)
	grams := ngram(input, lang, 2)
	print(grams)
	fmt.Println()
}

func tri(input string, lang alphabet.Language) {
	fmt.Printf("trigrams: %q\n", input)
	grams := ngram(input, lang, 3)
	print(grams)
	fmt.Println()
}

func quad(input string, lang alphabet.Language) {
	fmt.Printf("quadgrams: %q\n", input)
	grams := ngram(input, lang, 4)
	print(grams)
	fmt.Println()
}

func ngram(input string, lang alphabet.Language, size int) map[string]uint64 {
	result := make(map[string]uint64)
	buf := make([]rune, size)
	pos := 0
	count := 0

	rd := bufio.NewReader(strings.NewReader(input))

	for {
		r, _, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("err: %v\n", err)
			break
		}

		if unicode.IsSpace(r) {
			pos = 0
			count = 0
			continue
		}

		if !lang.ContainsRune(r) {
			continue
		}

		buf[pos+count] = r
		count++

		if count == size {
			gram := string(buf[pos:])

			copy(buf[pos:], buf[pos+1:])
			pos = 0
			count = size - 1

			freq, exists := result[gram]
			if !exists {
				result[gram] = 1
			} else {
				result[gram] = freq + 1
			}
		}
	}

	return result
}

func print(grams map[string]uint64) {
	pairs := collection.MapSortedByValue(grams, collection.Descending)

	for _, pair := range pairs {
		fmt.Printf("%s = %d\n", pair.Key, pair.Value)
	}
}
