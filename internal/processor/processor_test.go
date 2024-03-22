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
