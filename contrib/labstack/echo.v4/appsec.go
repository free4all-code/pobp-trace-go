

package echo

import (
	"net"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation/httpsec"

	"github.com/labstack/echo/v4"
)

func useAppSec(c echo.Context, span tracer.Span) func() {
	req := c.Request()
	httpsec.SetAppSecTags(span)
	params := make(map[string]string)
	for _, n := range c.ParamNames() {
		params[n] = c.Param(n)
	}
	args := httpsec.MakeHandlerOperationArgs(req, params)
	ctx, op := httpsec.StartOperation(req.Context(), args)
	c.SetRequest(req.WithContext(ctx))
	return func() {
		events := op.Finish(httpsec.HandlerOperationRes{Status: c.Response().Status})
		if len(events) > 0 {
			remoteIP, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				remoteIP = req.RemoteAddr
			}
			httpsec.SetSecurityEventTags(span, events, remoteIP, args.Headers, c.Response().Writer.Header())
		}
		instrumentation.SetTags(span, op.Tags())
	}
}
