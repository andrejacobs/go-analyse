package ngrams

import (
	"bufio"
	"context"
	"io"
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

	if tokenSize == 0 {
		return parseLetterMonograms(ctx, input, language, recv)
	}

	return parseLetterNgrams(ctx, input, language, tokenSize, recv)
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
