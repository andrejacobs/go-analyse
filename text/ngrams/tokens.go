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

package ngrams

import (
	"bufio"
	"context"
	"io"
	"strings"
	"unicode"

	"github.com/andrejacobs/go-analyse/text/alphabet"
)

// RecvTokenFunc will be called when a new token has been parsed from the input stream.
// The parser will pass any encountered error to this function and you should assume that the parsing process
// will stop and that no more tokens will be produced.
// If this function returns an error then it will indicate to the parser to stop the parsing process.
type RecvTokenFunc func(token string, err error) error

// ParseLetterTokens is used to parse ngrams for letter combinations of the given tokenSize and language from the io.Reader.
func ParseLetterTokens(ctx context.Context, input io.Reader, language alphabet.Language,
	tokenSize int, recv RecvTokenFunc) error {

	if tokenSize == 1 {
		return parseLetterMonograms(ctx, input, language, recv)
	}

	return parseLetterNgrams(ctx, input, language, tokenSize, recv)
}

// ParseWordTokens is used to parse ngrams for word combinations of the given tokenSize and language from the io.Reader.
func ParseWordTokens(ctx context.Context, input io.Reader, language alphabet.Language,
	tokenSize int, recv RecvTokenFunc) error {
	return parseWordNgrams(ctx, input, language, tokenSize, recv)
}

func parseLetterNgrams(ctx context.Context, input io.Reader, language alphabet.Language,
	tokenSize int, recv RecvTokenFunc) error {

	buf := make([]rune, tokenSize)
	pos := 0
	count := 0

	rd := bufio.NewReader(input)

loop:
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				// Inform consumer of error
				_ = recv("", err)
				return err
			}
			break loop
		default:
			r, _, err := rd.ReadRune()
			if err != nil {
				if err == io.EOF {
					// Done reading
					break loop
				}
				// Inform consumer of error
				_ = recv("", err)
				return err
			}

			// Ignore white space
			if unicode.IsSpace(r) {
				pos = 0
				count = 0
				continue
			}

			r = unicode.ToLower(r)

			// Ignore any runes not part of the language
			if !language.ContainsRune(r) {
				continue
			}

			buf[pos+count] = r
			count++

			// Did we parse enough runes for a full token?
			if count == tokenSize {
				token := string(buf[pos:])

				copy(buf[pos:], buf[pos+1:])
				pos = 0
				count = tokenSize - 1

				// Inform the consumer of a new token
				err := recv(token, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func parseLetterMonograms(ctx context.Context, input io.Reader,
	language alphabet.Language, recv RecvTokenFunc) error {

	rd := bufio.NewReader(input)

loop:
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				// Inform consumer of error
				_ = recv("", err)
				return err
			}
			break loop
		default:
			r, _, err := rd.ReadRune()
			if err != nil {
				if err == io.EOF {
					// Done reading
					break loop
				}
				// Inform consumer of error
				_ = recv("", err)
				return err
			}

			// Ignore white space
			if unicode.IsSpace(r) {
				continue
			}

			r = unicode.ToLower(r)

			// Ignore any runes not part of the language
			if !language.ContainsRune(r) {
				continue
			}

			// Inform the consumer of a new token
			err = recv(string(r), nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parseWordNgrams(ctx context.Context, input io.Reader, language alphabet.Language,
	tokenSize int, recv RecvTokenFunc) error {

	buf := make([]string, tokenSize)
	pos := 0
	count := 0

	scanner := bufio.NewScanner(bufio.NewReader(input))
	scanner.Split(bufio.ScanWords)

loop:
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				// Inform consumer of error
				_ = recv("", err)
				return err
			}
			break loop
		default:
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					// Inform consumer of error
					_ = recv("", err)
					return err
				}
				break loop
			}

			word := strings.ToLower(scanner.Text())

			buf[pos+count] = word
			count++

			// Did we parse enough words for a full token?
			if count == tokenSize {
				words := buf[pos:]
				token := strings.Join(words, " ")

				copy(buf[pos:], buf[pos+1:])
				pos = 0
				count = tokenSize - 1

				// Inform the consumer of a new token
				err := recv(token, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
