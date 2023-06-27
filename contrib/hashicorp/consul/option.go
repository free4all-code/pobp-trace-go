

package consul

import (
	"math"

	"git.proto.group/protoobp/pobp-trace-go/internal"
)

const (
	serviceName = "consul"
)

type clientConfig struct {
	serviceName   string
	analyticsRate float64
}

// ClientOption represents an option that can be used to create or wrap a client.
type ClientOption func(*clientConfig)

func defaults(cfg *clientConfig) {
	cfg.serviceName = serviceName
	if internal.BoolEnv("POBP_TRACE_CONSUL_ANALYTICS_ENABLED", false) {
		cfg.analyticsRate = 1.0
	} else {
		cfg.analyticsRate = math.NaN()
	}
}

// WithServiceName sets the given service name for the client.
func WithServiceName(name string) ClientOption {
	return func(cfg *clientConfig) {
		cfg.serviceName = name
	}
}

// WithAnalytics enables Trace Analytics for all started spans.
func WithAnalytics(on bool) ClientOption {
	return func(cfg *clientConfig) {
		if on {
			cfg.analyticsRate = 1.0
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

// WithAnalyticsRate sets the sampling rate for Trace Analytics events
// correlated to started spans.
func WithAnalyticsRate(rate float64) ClientOption {
	return func(cfg *clientConfig) {
		if rate >= 0.0 && rate <= 1.0 {
			cfg.analyticsRate = rate
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}