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

// Package compiledinfo is an internal package used to bake "build time information" into the compiled app.
package compiledinfo

import "fmt"

const (
	defaultAppName = "TODO"
)

var (
	// The name of the app (binary executable).
	AppName string
	// The hash of the last commit in the git repository.
	GitCommitHash string
	// The version of the app.
	Version string
)

// Return the version information that can be shown to a user.
func VersionString() string {
	version := "v0.0.0"
	if Version != "" {
		version = Version
	}
	return fmt.Sprintf("%s %s", version, GitCommitHash)
}

// Return the name of the app as shown in help and version output.
func UsageName() string {
	appName := defaultAppName
	if AppName != "" {
		appName = AppName
	}
	return appName
}

// Return the name of the app and version info as displayed in the usage information.
func UsageNameAndVersion() string {
	return fmt.Sprintf("%s version: %s", UsageName(), VersionString())
}
