

//go:build !windows && !linux && !darwin && !freebsd
// +build !windows,!linux,!darwin,!freebsd

package osinfo

import (
	"runtime"
)

func osName() string {
	return runtime.GOOS
}

func osVersion() string {
	return "unknown"
}
