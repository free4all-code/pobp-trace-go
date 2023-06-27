

package restful

import (
	"math"

	"git.proto.group/protoobp/pobp-trace-go/internal"
	"git.proto.group/protoobp/pobp-trace-go/internal/globalconfig"
)

type config struct {
	serviceName   string
	analyticsRate float64
}

func newConfig() *config {
	rate := globalconfig.AnalyticsRate()
	if internal.BoolEnv("POBP_TRACE_RESTFUL_ANALYTICS_ENABLED", false) {
		rate = 1.0
	}
	return &config{
		serviceName:   "go-restful",
		analyticsRate: rate,
	}
}

// Option specifies instrumentation configuration options.
type Option func(*config)

// WithServiceName sets the service name to by used by the filter.
func WithServiceName(name string) Option {
	return func(cfg *config) {
		cfg.serviceName = name
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
