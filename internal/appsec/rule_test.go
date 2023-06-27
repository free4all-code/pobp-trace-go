

//go:build appsec
// +build appsec

package appsec

import (
	"testing"

	"github.com/stretchr/testify/require"

	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/waf"
)

func TestStaticRule(t *testing.T) {
	if waf.Health() != nil {
		t.Skip("waf disabled")
		return
	}
	waf, err := waf.NewHandle([]byte(staticRecommendedRule), "", "")
	require.NoError(t, err)
	waf.Close()
}
