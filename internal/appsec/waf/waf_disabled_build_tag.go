

// Build when CGO is enabled but the target OS or architecture are not supported
//go:build !appsec
// +build !appsec

package waf

var disabledReason = "the waf is disabled due to missing go build tag appsec"
