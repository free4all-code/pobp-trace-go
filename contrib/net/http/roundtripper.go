

package http

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
)

type roundTripper struct {
	base http.RoundTripper
	cfg  *roundTripperConfig
}

func (rt *roundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	resourceName := rt.cfg.resourceNamer(req)
	opts := []pobptrace.StartSpanOption{
		tracer.SpanType(ext.SpanTypeHTTP),
		tracer.ResourceName(resourceName),
		tracer.Tag(ext.HTTPMethod, req.Method),
		tracer.Tag(ext.HTTPURL, req.URL.String()),
	}
	if !math.IsNaN(rt.cfg.analyticsRate) {
		opts = append(opts, tracer.Tag(ext.EventSampleRate, rt.cfg.analyticsRate))
	}
	if rt.cfg.serviceName != "" {
		opts = append(opts, tracer.ServiceName(rt.cfg.serviceName))
	}
	if len(rt.cfg.spanOpts) > 0 {
		opts = append(opts, rt.cfg.spanOpts...)
	}
	span, ctx := tracer.StartSpanFromContext(req.Context(), "http.request", opts...)
	defer func() {
		if rt.cfg.after != nil {
			rt.cfg.after(res, span)
		}
		span.Finish(tracer.WithError(err))
	}()
	if rt.cfg.before != nil {
		rt.cfg.before(req, span)
	}
	r2 := req.Clone(ctx)
	// inject the span context into the http request copy
	err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(r2.Header))
	if err != nil {
		// this should never happen
		fmt.Fprintf(os.Stderr, "contrib/net/http.Roundtrip: failed to inject http headers: %v\n", err)
	}
	res, err = rt.base.RoundTrip(r2)
	if err != nil {
		span.SetTag("http.errors", err.Error())
		span.SetTag(ext.Error, err)
	} else {
		span.SetTag(ext.HTTPCode, strconv.Itoa(res.StatusCode))
		// treat 5XX as errors
		if res.StatusCode/100 == 5 {
			span.SetTag("http.errors", res.Status)
			span.SetTag(ext.Error, fmt.Errorf("%d: %s", res.StatusCode, http.StatusText(res.StatusCode)))
		}
	}
	return res, err
}

// Unwrap returns the original http.RoundTripper.
func (rt *roundTripper) Unwrap() http.RoundTripper {
	return rt.base
}

// WrapRoundTripper returns a new RoundTripper which traces all requests sent
// over the transport.
func WrapRoundTripper(rt http.RoundTripper, opts ...RoundTripperOption) http.RoundTripper {
	cfg := newRoundTripperConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	if wrapped, ok := rt.(*roundTripper); ok {
		rt = wrapped.base
	}
	return &roundTripper{
		base: rt,
		cfg:  cfg,
	}
}

// WrapClient modifies the given client's transport to augment it with tracing and returns it.
func WrapClient(c *http.Client, opts ...RoundTripperOption) *http.Client {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
	c.Transport = WrapRoundTripper(c.Transport, opts...)
	return c
}