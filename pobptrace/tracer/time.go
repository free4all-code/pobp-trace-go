

//go:build !windows
// +build !windows

package tracer

import "time"

// now returns the current UNIX time in nanoseconds, as computed by Time.UnixNano().
func now() int64 {
	return time.Now().UnixNano()
}
