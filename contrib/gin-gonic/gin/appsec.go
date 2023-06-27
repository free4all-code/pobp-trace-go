
package gin

import (
	"net"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation/httpsec"

	"github.com/gin-gonic/gin"
)

// useAppSec executes the AppSec logic related to the operation start and
// returns the  function to be executed upon finishing the operation
func useAppSec(c *gin.Context, span tracer.Span) func() {
	req := c.Request
	httpsec.SetAppSecTags(span)
	var params map[string]string
	if l := len(c.Params); l > 0 {
		params = make(map[string]string, l)
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}
	}
	args := httpsec.MakeHandlerOperationArgs(req, params)
	ctx, op := httpsec.StartOperation(req.Context(), args)
	c.Request = req.WithContext(ctx)
	return func() {
		events := op.Finish(httpsec.HandlerOperationRes{Status: c.Writer.Status()})
		if len(events) > 0 {
			remoteIP, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				remoteIP = req.RemoteAddr
			}
			httpsec.SetSecurityEventTags(span, events, remoteIP, args.Headers, c.Writer.Header())
		}
		instrumentation.SetTags(span, op.Tags())
	}
}
