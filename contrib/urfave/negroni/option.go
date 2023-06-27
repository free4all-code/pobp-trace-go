

// Package negroni provides helper functions for tracing the urfave/negroni package (https://github.com/urfave/negroni).
package negroni

import (
	"math"
	"net/http"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/internal"
	"git.proto.group/protoobp/pobp-trace-go/internal/globalconfig"
)

type config struct {
	serviceName   string
	spanOpts      []pobptrace.StartSpanOption // additional span options to be applied
	analyticsRate float64
	isStatusError func(statusCode int) bool
	resourceNamer func(r *http.Request) string
}

// Option represents an option that can be passed to NewRouter.
type Option func(*config)

func defaults(cfg *config) {
	cfg.serviceName = "negroni.router"
	if svc := globalconfig.ServiceName(); svc != "" {
		cfg.serviceName = svc
	}
	if internal.BoolEnv("POBP_TRACE_NEGRONI_ANALYTICS_ENABLED", false) {
		cfg.analyticsRate = 1.0
	} else {
		cfg.analyticsRate = globalconfig.AnalyticsRate()
	}
	cfg.isStatusError = isServerError
	cfg.resourceNamer = defaultResourceNamer
}

// WithServiceName sets the given service name for the router.
func WithServiceName(name string) Option {
	return func(cfg *config) {
		cfg.serviceName = name
	}
}

// WithSpanOptions applies the given set of options to the spans started
// by the router.
func WithSpanOptions(opts ...pobptrace.StartSpanOption) Option {
	return func(cfg *config) {
		cfg.spanOpts = opts
	}
}

// WithAnalytics enables Trace Analytics for all started spans.
func WithAnalytics(on bool) Option {
	return func(cfg *config) {
		if on {
			cfg.analyticsRate = 1.0
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

// WithAnalyticsRate sets the sampling rate for Trace Analytics events
// correlated to started spans.
func WithAnalyticsRate(rate float64) Option {
	return func(cfg *config) {
		if rate >= 0.0 && rate <= 1.0 {
			cfg.analyticsRate = rate
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

// WithStatusCheck specifies a function fn which reports whether the passed
// statusCode should be considered an error.
func WithStatusCheck(fn func(statusCode int) bool) Option {
	return func(cfg *config) {
		cfg.isStatusError = fn
	}
}

func isServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

// WithResourceNamer specifies a function which will be used to obtain a resource name for a given
// negroni request, using the request's context.
func WithResourceNamer(namer func(r *http.Request) string) Option {
	return func(cfg *config) {
		cfg.resourceNamer = namer
	}
}

func defaultResourceNamer(r *http.Request) string {
	return ""
}