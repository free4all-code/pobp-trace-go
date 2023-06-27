

// Package http provides functions to trace the net/http package (https://golang.org/pkg/net/http).
package http // import "git.proto.group/protoobp/pobp-trace-go/contrib/net/http"

import (
	"net/http"

	"git.proto.group/protoobp/pobp-trace-go/internal/log"
)

// ServeMux is an HTTP request multiplexer that traces all the incoming requests.
type ServeMux struct {
	*http.ServeMux
	cfg *config
}

// NewServeMux allocates and returns an http.ServeMux augmented with the
// global tracer.
func NewServeMux(opts ...Option) *ServeMux {
	cfg := new(config)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	log.Debug("contrib/net/http: Configuring ServeMux: %#v", cfg)
	return &ServeMux{
		ServeMux: http.NewServeMux(),
		cfg:      cfg,
	}
}

// ServeHTTP dispatches the request to the handler
// whose pattern most closely matches the request URL.
// We only need to rewrite this function to be able to trace
// all the incoming requests to the underlying multiplexer
func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mux.cfg.ignoreRequest(r) {
		mux.ServeMux.ServeHTTP(w, r)
		return
	}
	// get the resource associated to this request
	_, route := mux.Handler(r)
	resource := r.Method + " " + route
	TraceAndServe(mux.ServeMux, w, r, &ServeConfig{
		Service:  mux.cfg.serviceName,
		Resource: resource,
		SpanOpts: mux.cfg.spanOpts,
		Route:    route,
	})
}

// WrapHandler wraps an http.Handler with tracing using the given service and resource.
// If the WithResourceNamer option is provided as part of opts, it will take precedence over the resource argument.
func WrapHandler(h http.Handler, service, resource string, opts ...Option) http.Handler {
	cfg := new(config)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	log.Debug("contrib/net/http: Wrapping Handler: Service: %s, Resource: %s, %#v", service, resource, cfg)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if cfg.ignoreRequest(req) {
			h.ServeHTTP(w, req)
			return
		}
		if r := cfg.resourceNamer(req); r != "" {
			resource = r
		}
		TraceAndServe(h, w, req, &ServeConfig{
			Service:    service,
			Resource:   resource,
			FinishOpts: cfg.finishOpts,
			SpanOpts:   cfg.spanOpts,
			Route:      req.URL.EscapedPath(),
		})
	})
}
