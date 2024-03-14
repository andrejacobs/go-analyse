package alphabet_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/stretchr/testify/assert"
)

func TestDiscoverLetters(t *testing.T) {
	r := strings.NewReader(`
	Goeie môre, my vrou,
	The rabbit-hole went straight on like a tunnel for some way, and then
	qu'elles étaient garnies d'armoires et d'étagères; çà et là, elle vit
	武士道改善
`)

	expected := []rune("',-;abdefghiklmnoqrstuvwyàçèéô善士改武道")

	letters, err := alphabet.DiscoverLetters(context.Background(), r)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, letters)
}

func TestDiscoverLettersContextCancelled(t *testing.T) {
	r := strings.NewReader(`The quick brown fox jumped over the lazy dog!`)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := alphabet.DiscoverLetters(ctx, r)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestDiscoverLettersReadError(t *testing.T) {
	var r FailReader
	_, err := alphabet.DiscoverLetters(context.Background(), &r)
	assert.ErrorContains(t, err, "failed to read")
}

func TestDiscoverLettersFromFile(t *testing.T) {
	expected := []rune("',-;abdefghiklmnoqrstuvwyàçèéô善士改武道")

	letters, err := alphabet.DiscoverLettersFromFile(context.Background(), "testdata/discover.txt")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, letters)

	_, err = alphabet.DiscoverLettersFromFile(context.Background(), "testdata/naf.txt")
	assert.ErrorContains(t, err, "failed to open \"testdata/naf.txt\"")
}

// -----------------------------------------------------------------------------
type FailReader bool

func (fr *FailReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to read")
}
