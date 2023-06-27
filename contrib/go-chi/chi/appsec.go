

package chi

import (
	"net/http"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation/httpsec"

	"github.com/go-chi/chi"
)

func withAppsec(next http.Handler, r *http.Request, span tracer.Span) http.Handler {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return httpsec.WrapHandler(next, span, nil)
	}
	var pathParams map[string]string
	keys := rctx.URLParams.Keys
	values := rctx.URLParams.Values
	if len(keys) > 0 && len(keys) == len(values) {
		pathParams = make(map[string]string, len(keys))
		for i, key := range keys {
			pathParams[key] = values[i]
		}
	}
	return httpsec.WrapHandler(next, span, pathParams)
}
