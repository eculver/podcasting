// The version package provides a location to set the release versions for all
// packages to consume, without creating import cycles.
//
// This package should not import any other tdtv2 packages.
//
// TODO: replace this with native binary versioning support in Go 1.18
package version

import (
	version "github.com/hashicorp/go-version"
)

// Version is The main version number that is being run at the moment.
var Version = "0.2.0"

// SemVer is an instance of version.Version. This has the secondary
// benefit of verifying during tests and init time that our version is a
// proper semantic version, which should always be the case.
var SemVer *version.Version

func init() {
	SemVer = version.Must(version.NewVersion(Version))
}

// String returns the complete version string
func String() string {
	return SemVer.String()
}
