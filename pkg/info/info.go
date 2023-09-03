/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// This file contains information about the tool.

package info

import (
	"fmt"
	"runtime"
)

// Version represents the version information of the tool, including
// the major, minor, and patch versions, as well as the Go runtime version.
type Version struct {
	Major     string
	Minor     string
	Patch     string
	GoVersion string
}

// GetGoVersion returns the current Go runtime version as a string.
func GetGoVersion() string {
	return runtime.Version()
}

// BuildCurrentVersion builds and returns the current Version of the tool.
func BuildCurrentVersion() Version {
	v := Version{
		Major:     "1",
		Minor:     "0",
		Patch:     "0",
		GoVersion: GetGoVersion(),
	}
	return v
}

// PrintVersion prints the current version of the tool to standard output.
func PrintVersion() {
	v := BuildCurrentVersion()
	fmt.Printf("Version: {Major:\"%s\", Minor:\"%s\", Patch:\"%s\", GoVersion:\"%s\"}\n", v.Major, v.Minor, v.Patch, v.GoVersion)
}
