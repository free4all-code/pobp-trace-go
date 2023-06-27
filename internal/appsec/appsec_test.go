

//go:build appsec
// +build appsec

package appsec_test

import (
	"os"
	"strconv"
	"testing"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/waf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnabled(t *testing.T) {
	enabledConfig, _ := strconv.ParseBool(os.Getenv("POBP_APPSEC_ENABLED"))
	canBeEnabled := enabledConfig && waf.Health() == nil

	require.False(t, appsec.Enabled())
	tracer.Start()
	assert.Equal(t, canBeEnabled, appsec.Enabled())
	tracer.Stop()
	assert.False(t, appsec.Enabled())
}
