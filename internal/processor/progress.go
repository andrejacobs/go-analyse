package processor

import "io"

// ProgressReporter is used to report progress updates.
type ProgressReporter interface {
	// Started will be called when a new path is being processed.
	// index is the 0th based index of the path in the total number of paths.
	// path uniquely identifies a resource and does not have to be a file path (e.g. a URI)
	Started(path string, index int, total int)

	// Reader returns a new wrapped reader that will update and report progress as data
	// is being read from it.
	Reader(r io.Reader) io.Reader

	// AddToTotalSize is called when the total number of bytes to be processed has changed.
	// For example like reading from a zip file.
	AddToTotalSize(add int64)
}

//-----------------------------------------------------------------------------

// ProgressReporter implementation that does nothing.
type nullProgressReporter struct {
}

func (n *nullProgressReporter) Started(path string, index int, total int) {
}

func (n *nullProgressReporter) Reader(r io.Reader) io.Reader {
	return r
}

func (n *nullProgressReporter) AddToTotalSize(add int64) {
}
