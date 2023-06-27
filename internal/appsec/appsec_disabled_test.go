

//go:build !appsec
// +build !appsec

package appsec_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec"

	"github.com/stretchr/testify/require"
)

func TestEnabled(t *testing.T) {
	enabledStr := os.Getenv("POBP_APPSEC_ENABLED")
	if enabledStr != "" {
		defer os.Setenv("POBP_APPSEC_ENABLED", enabledStr)
	}
	// AppSec should be always disabled
	require.False(t, appsec.Enabled())
	tracer.Start()
	assert.False(t, appsec.Enabled())
	tracer.Stop()
	assert.False(t, appsec.Enabled())
	os.Setenv("POBP_APPSEC_ENABLED", "true")
	require.False(t, appsec.Enabled())
	tracer.Start()
	assert.False(t, appsec.Enabled())
	tracer.Stop()
	assert.False(t, appsec.Enabled())

}
