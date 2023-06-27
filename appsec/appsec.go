

package appsec

import (
	"context"

	"git.proto.group/protoobp/pobp-trace-go/internal/appsec"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation/httpsec"
)

// MonitorParsedHTTPBody runs the security monitoring rules on the given *parsed*
// HTTP request body. The given context must be the HTTP request context as returned
// by the Context() method of an HTTP request. Calls to this function are ignored if
// AppSec is disabled or the given context is incorrect.
// Note that passing the raw bytes of the HTTP request body is not expected and would
// result in inaccurate attack detection.
func MonitorParsedHTTPBody(ctx context.Context, body interface{}) {
	if appsec.Enabled() {
		httpsec.MonitorParsedBody(ctx, body)
	}
	// bonus: use sync.Once to log a debug message once if AppSec is disabled
}
