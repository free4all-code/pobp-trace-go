

// Build when CGO is enabled but the target OS or architecture are not supported
//go:build appsec && cgo && (windows || !amd64)
// +build appsec
// +build cgo
// +build windows !amd64

package waf

import (
	"fmt"
	"runtime"
)

var disabledReason = fmt.Sprintf("the target operating-system %s or architecture %s are not supported", runtime.GOOS, runtime.GOARCH)
