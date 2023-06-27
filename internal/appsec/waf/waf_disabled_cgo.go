

// Build when CGO is disabled
//go:build appsec && !cgo
// +build appsec,!cgo

package waf

var disabledReason = "cgo was disabled during the compilation and should be enabled in order to compile with the waf"
