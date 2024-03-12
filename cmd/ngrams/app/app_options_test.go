package app

import (
	"testing"

	"github.com/andrejacobs/go-analyse/text/alphabet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionsWithDefaults(t *testing.T) {
	var opt options
	require.NoError(t, WithDefaults()(&opt))

	assert.Equal(t, alphabet.LanguageCode("en"), opt.langCode)
	assert.False(t, opt.words)
	assert.Equal(t, 1, opt.tokenSize)
}

func TestOptionsWithLettersOrWords(t *testing.T) {
	var opt options
	require.NoError(t, WithLetters()(&opt))
	assert.False(t, opt.words)

	require.NoError(t, WithWords()(&opt))
	assert.True(t, opt.words)
}
