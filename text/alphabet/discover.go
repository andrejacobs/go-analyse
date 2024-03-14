package alphabet

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"unicode"

	"github.com/andrejacobs/go-collection/collection"
)

// DiscoverLetters produces a slice containing the unique non-whitespace lowercased letters found in the io.Reader.
func DiscoverLetters(ctx context.Context, input io.Reader) ([]rune, error) {

	result := collection.NewSet[rune]()
	rd := bufio.NewReader(input)

loop:
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return result.Items(), err
			}
			break loop
		default:
			r, _, err := rd.ReadRune()
			if err != nil {
				if err == io.EOF {
					// Done reading
					break loop
				}
				return result.Items(), err
			}

			// Ignore white space
			if unicode.IsSpace(r) {
				continue
			}

			result.Insert(unicode.ToLower(r))
		}
	}

	return result.Items(), nil
}

// DiscoverLettersFromFile produces a slice containing the unique non-whitespace lowercased letters found in the file.
func DiscoverLettersFromFile(ctx context.Context, path string) ([]rune, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q. %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to close %q. %v", path, err)
		}
	}()

	result, err := DiscoverLetters(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to discover letters from %q. %w", path, err)
	}

	return result, nil
}
