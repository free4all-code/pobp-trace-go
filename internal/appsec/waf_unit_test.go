

//go:build appsec
// +build appsec

package appsec

import (
	"testing"

	"github.com/stretchr/testify/require"

	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/waf"
)

// Test that internal functions used to set span tags use the correct types
func TestTagsTypes(t *testing.T) {
	th := instrumentation.NewTagsHolder()
	rInfo := waf.RulesetInfo{
		Version: "1.3.0",
		Loaded:  10,
		Failed:  1,
		Errors:  map[string]interface{}{"test": []string{"1", "2"}},
	}

	addRulesMonitoringTags(&th, rInfo)
	addWAFMonitoringTags(&th, "1.2.3", 2, 1, 3)

	tags := th.Tags()
	_, ok := tags[eventRulesErrorsTag].(string)
	require.True(t, ok)

	for _, tag := range []string{eventRulesLoadedTag, eventRulesFailedTag, wafDurationTag, wafDurationExtTag, wafVersionTag} {
		require.Contains(t, tags, tag)
	}
}
