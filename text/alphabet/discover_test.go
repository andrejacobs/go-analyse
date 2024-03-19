package alphabet_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestDiscoverProcessor(t *testing.T) {
	expected := []rune("',-;abdefghiklmnoqrstuvwyàçèéô善士改武道")

	p := alphabet.NewDiscoverProcessor()
	err := p.ProcessFiles(context.Background(), []string{"testdata/discover.txt"})
	require.NoError(t, err)
	assert.ElementsMatch(t, expected, p.Letters())

	temp := filepath.Join(os.TempDir(), "alphabet-unit-test-lang.csv")
	defer os.Remove(temp)
	require.NoError(t, p.Save(temp))

	am, err := alphabet.LoadLanguagesFromFile(temp)
	require.NoError(t, err)
	lang, err := am.Get("unknown")
	assert.NoError(t, err)
	assert.Equal(t, "unknown", lang.Name)
	assert.Equal(t, alphabet.LanguageCode("unknown"), lang.Code)
	assert.ElementsMatch(t, expected, []rune(lang.Letters))
}

// -----------------------------------------------------------------------------
type FailReader bool

func (fr *FailReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to read")
}
