

// Build when CGO is disabled or the target OS or Arch are not supported
//go:build !appsec || !cgo || windows || !amd64
// +build !appsec !cgo windows !amd64

package waf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	require.Error(t, Health())
}
