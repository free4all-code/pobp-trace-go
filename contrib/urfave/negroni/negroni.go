
// Package negroni provides helper functions for tracing the urfave/negroni package (https://github.com/urfave/negroni).
package negroni

import (
	"fmt"
	"math"
	"net/http"

	"github.com/urfave/negroni"

	"git.proto.group/protoobp/pobp-trace-go/contrib/internal/httptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/log"
)

type DatadogMiddleware struct {
	cfg *config
}

func (m *DatadogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	opts := append(m.cfg.spanOpts, tracer.ServiceName(m.cfg.serviceName), tracer.ResourceName(m.cfg.resourceNamer(r)))
	if !math.IsNaN(m.cfg.analyticsRate) {
		opts = append(opts, tracer.Tag(ext.EventSampleRate, m.cfg.analyticsRate))
	}
	span, ctx := httptrace.StartRequestSpan(r, opts...)
	defer func() {
		// check if the responseWriter is of type negroni.ResponseWriter
		var (
			status int
			opts   []tracer.FinishOption
		)
		responseWriter, ok := w.(negroni.ResponseWriter)
		if ok {
			status = responseWriter.Status()
			if m.cfg.isStatusError(status) {
				opts = []tracer.FinishOption{tracer.WithError(fmt.Errorf("%d: %s", status, http.StatusText(status)))}
			}
		}
		httptrace.FinishRequestSpan(span, status, opts...)
	}()

	next(w, r.WithContext(ctx))
}

// Middleware create the negroni middleware that will trace incoming requests
func Middleware(opts ...Option) *DatadogMiddleware {
	cfg := new(config)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	log.Debug("contrib/urgave/negroni: Configuring Middleware: %#v", cfg)

	m := DatadogMiddleware{
		cfg: cfg,
	}

	return &m
}
