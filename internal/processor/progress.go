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
