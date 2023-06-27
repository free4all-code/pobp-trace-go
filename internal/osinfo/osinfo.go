
// Package osinfo provides information about the current operating system release
package osinfo

// OSName returns the name of the operating system, including the distribution
// for Linux when possible.
func OSName() string {
	// call out to OS-specific implementation
	return osName()
}

// OSVersion returns the operating system release, e.g. major/minor version
// number and build ID.
func OSVersion() string {
	// call out to OS-specific implementation
	return osVersion()
}
