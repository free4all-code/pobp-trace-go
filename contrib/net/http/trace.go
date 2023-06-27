

package http // import "git.proto.group/protoobp/pobp-trace-go/contrib/net/http"

//go:generate sh -c "go run make_responsewriter.go | gofmt > trace_gen.go"

import (
	"net/http"

	"git.proto.group/protoobp/pobp-trace-go/contrib/internal/httptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation/httpsec"
)

// ServeConfig specifies the tracing configuration when using TraceAndServe.
type ServeConfig struct {
	// Service specifies the service name to use. If left blank, the global service name
	// will be inherited.
	Service string
	// Resource optionally specifies the resource name for this request.
	Resource string
	// QueryParams should be true in order to append the URL query values to the  "http.url" tag.
	QueryParams bool
	// Route is the request matched route if any, or is empty otherwise
	Route string
	// RouteParams specifies framework-specific route parameters (e.g. for route /user/:id coming
	// in as /user/123 we'll have {"id": "123"}). This field is optional and is used for monitoring
	// by AppSec. It is only taken into account when AppSec is enabled.
	RouteParams map[string]string
	// FinishOpts specifies any options to be used when finishing the request span.
	FinishOpts []pobptrace.FinishOption
	// SpanOpts specifies any options to be applied to the request starting span.
	SpanOpts []pobptrace.StartSpanOption
}

// TraceAndServe serves the handler h using the given ResponseWriter and Request, applying tracing
// according to the specified config.
func TraceAndServe(h http.Handler, w http.ResponseWriter, r *http.Request, cfg *ServeConfig) {
	if cfg == nil {
		cfg = new(ServeConfig)
	}
	opts := append(cfg.SpanOpts, tracer.ServiceName(cfg.Service), tracer.ResourceName(cfg.Resource))
	opts = append(opts, tracer.Tag(ext.HTTPRoute, cfg.Route))
	span, ctx := httptrace.StartRequestSpan(r, opts...)
	rw, ddrw := wrapResponseWriter(w)
	defer func() {
		httptrace.FinishRequestSpan(span, ddrw.status, cfg.FinishOpts...)
	}()

	if appsec.Enabled() {
		h = httpsec.WrapHandler(h, span, cfg.RouteParams)
	}
	h.ServeHTTP(rw, r.WithContext(ctx))
}

// responseWriter is a small wrapper around an http response writer that will
// intercept and store the status of a request.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, 0}
}

// Status returns the status code that was monitored.
func (w *responseWriter) Status() int {
	return w.status
}

// Write writes the data to the connection as part of an HTTP reply.
// We explicitly call WriteHeader with the 200 status code
// in order to get it reported into the span.
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

// WriteHeader sends an HTTP response header with status code.
// It also sets the status code to the span.
func (w *responseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}