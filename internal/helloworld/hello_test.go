package helloworld_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/andrejacobs/go-analysis/internal/helloworld"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSayHello(t *testing.T) {
	name := "Andre"
	expected := fmt.Sprintf("Hello %s, it is so nice to meet you.\n", name)

	var buf bytes.Buffer

	require.NoError(t, helloworld.SayHello(&buf, name))
	assert.Equal(t, expected, buf.String())
}
