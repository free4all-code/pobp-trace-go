

// Package mux provides tracing functions for tracing the gorilla/mux package (https://github.com/gorilla/mux).
package mux // import "git.proto.group/protoobp/pobp-trace-go/contrib/gorilla/mux"

import (
	"net/http"
	"strings"

	httptrace "git.proto.group/protoobp/pobp-trace-go/contrib/net/http"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/log"

	"github.com/gorilla/mux"
)

// Router registers routes to be matched and dispatches a handler.
type Router struct {
	*mux.Router
	config *routerConfig
}

// StrictSlash defines the trailing slash behavior for new routes. The initial
// value is false.
//
// When true, if the route path is "/path/", accessing "/path" will perform a redirect
// to the former and vice versa. In other words, your application will always
// see the path as specified in the route.
//
// When false, if the route path is "/path", accessing "/path/" will not match
// this route and vice versa.
//
// The re-direct is a HTTP 301 (Moved Permanently). Note that when this is set for
// routes with a non-idempotent method (e.g. POST, PUT), the subsequent re-directed
// request will be made as a GET by most clients. Use middleware or client settings
// to modify this behaviour as needed.
//
// Special case: when a route sets a path prefix using the PathPrefix() method,
// strict slash is ignored for that route because the redirect behavior can't
// be determined from a prefix alone. However, any subrouters created from that
// route inherit the original StrictSlash setting.
func (r *Router) StrictSlash(value bool) *Router {
	r.Router.StrictSlash(value)
	return r
}

// SkipClean defines the path cleaning behaviour for new routes. The initial
// value is false. Users should be careful about which routes are not cleaned
//
// When true, if the route path is "/path//to", it will remain with the double
// slash. This is helpful if you have a route like: /fetch/http://xkcd.com/534/
//
// When false, the path will be cleaned, so /fetch/http://xkcd.com/534/ will
// become /fetch/http/xkcd.com/534
func (r *Router) SkipClean(value bool) *Router {
	r.Router.SkipClean(value)
	return r
}

// UseEncodedPath tells the router to match the encoded original path
// to the routes.
// For eg. "/path/foo%2Fbar/to" will match the path "/path/{var}/to".
//
// If not called, the router will match the unencoded path to the routes.
// For eg. "/path/foo%2Fbar/to" will match the path "/path/foo/bar/to"
func (r *Router) UseEncodedPath() *Router {
	r.Router.UseEncodedPath()
	return r
}

// NewRouter returns a new router instance traced with the global tracer.
func NewRouter(opts ...RouterOption) *Router {
	return WrapRouter(mux.NewRouter(), opts...)
}

// ServeHTTP dispatches the request to the handler
// whose pattern most closely matches the request URL.
// We only need to rewrite this function to be able to trace
// all the incoming requests to the underlying multiplexer
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.config.ignoreRequest(req) {
		r.Router.ServeHTTP(w, req)
		return
	}
	var (
		match    mux.RouteMatch
		spanopts []pobptrace.StartSpanOption
		route    string
	)
	// get the resource associated to this request
	if r.Match(req, &match) && match.Route != nil {
		if h, err := match.Route.GetHostTemplate(); err == nil {
			spanopts = append(spanopts, tracer.Tag("mux.host", h))
		}
		route, _ = match.Route.GetPathTemplate()
	}
	spanopts = append(spanopts, r.config.spanOpts...)
	if r.config.headerTags {
		spanopts = append(spanopts, headerTagsFromRequest(req))
	}
	resource := r.config.resourceNamer(r, req)
	httptrace.TraceAndServe(r.Router, w, req, &httptrace.ServeConfig{
		Service:     r.config.serviceName,
		Resource:    resource,
		FinishOpts:  r.config.finishOpts,
		SpanOpts:    spanopts,
		QueryParams: r.config.queryParams,
		RouteParams: match.Vars,
		Route:       route,
	})
}

// WrapRouter returns the given router wrapped with the tracing of the HTTP
// requests and responses served by the router.
func WrapRouter(router *mux.Router, opts ...RouterOption) *Router {
	cfg := newConfig(opts)
	log.Debug("contrib/gorilla/mux: Configuring Router: %#v", cfg)
	return &Router{
		Router: router,
		config: cfg,
	}
}

// defaultResourceNamer attempts to quantize the resource for an HTTP request by
// retrieving the path template associated with the route from the request.
func defaultResourceNamer(router *Router, req *http.Request) string {
	var match mux.RouteMatch
	// get the resource associated with the given request
	if router.Match(req, &match) && match.Route != nil {
		if r, err := match.Route.GetPathTemplate(); err == nil {
			return req.Method + " " + r
		}
	}
	return req.Method + " unknown"
}

func headerTagsFromRequest(req *http.Request) pobptrace.StartSpanOption {
	return func(cfg *pobptrace.StartSpanConfig) {
		for k := range req.Header {
			if !strings.HasPrefix(strings.ToLower(k), "x-protoobp-") {
				cfg.Tags["http.request.headers."+k] = strings.Join(req.Header.Values(k), ",")
			}
		}
	}
}