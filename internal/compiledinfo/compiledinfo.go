// Package compiledinfo is an internal package used to bake "build time information" into the compiled app.
package compiledinfo

import "fmt"

const (
	defaultAppName = "helloworld"
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
