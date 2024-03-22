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

package processor_test

import (
	"context"
	"io"
	"testing"

	"github.com/andrejacobs/go-analyse/internal/processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessor(t *testing.T) {
	paths := []string{"testdata/1.txt", "testdata/2.txt"}

	result := ""
	p := processor.NewProcessor()
	err := p.ProcessFiles(context.Background(), paths, func(ctx context.Context, r io.Reader) error {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		result += string(data)
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, "The quick brown foxjumped over the lazy dog!", result)
}

func TestProcessorWithProgress(t *testing.T) {
	paths := []string{"testdata/1.txt", "testdata/2.txt"}

	reporter := MockProgressReporter{}
	result := ""
	p := processor.NewProcessor()
	p.SetProgressReporter(&reporter)
	err := p.ProcessFiles(context.Background(), paths, func(ctx context.Context, r io.Reader) error {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		result += string(data)
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, "The quick brown foxjumped over the lazy dog!", result)
	assert.Equal(t, 2, reporter.startedTotal)
	assert.Equal(t, 2, reporter.startedCalled)
	assert.True(t, reporter.readerCalled)
	assert.Equal(t, int64(19+25), reporter.addTotal)
}

func TestProcessorZip(t *testing.T) {
	paths := []string{"testdata/a.zip"}

	result := ""
	p := processor.NewProcessor()
	err := p.ProcessFiles(context.Background(), paths, func(ctx context.Context, r io.Reader) error {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		result += string(data)
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, "The quick brown foxjumped over the lazy dog!", result)
}

func TestProcessorZipWithProgress(t *testing.T) {
	paths := []string{"testdata/a.zip"}

	reporter := MockProgressReporter{}
	result := ""
	p := processor.NewProcessor()
	p.SetProgressReporter(&reporter)
	err := p.ProcessFiles(context.Background(), paths, func(ctx context.Context, r io.Reader) error {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		result += string(data)
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, "The quick brown foxjumped over the lazy dog!", result)
	assert.Equal(t, 1, reporter.startedTotal)
	assert.Equal(t, 1, reporter.startedCalled)
	assert.True(t, reporter.readerCalled)
	assert.Equal(t, int64(19+25), reporter.addTotal)
}

//-----------------------------------------------------------------------------

type MockProgressReporter struct {
	startedTotal  int
	startedCalled int
	readerCalled  bool
	addTotal      int64
}

func (n *MockProgressReporter) Started(path string, index int, total int) {
	n.startedTotal = total
	n.startedCalled++
}

func (n *MockProgressReporter) Reader(r io.Reader) io.Reader {
	n.readerCalled = true
	return r
}

func (n *MockProgressReporter) AddToTotalSize(add int64) {
	n.addTotal += add
}
