
package redis

import (
	"math"

	"git.proto.group/protoobp/pobp-trace-go/internal"
)

type clientConfig struct {
	serviceName   string
	analyticsRate float64
	skipRaw       bool
}

// ClientOption represents an option that can be used to create or wrap a client.
type ClientOption func(*clientConfig)

func defaults(cfg *clientConfig) {
	cfg.serviceName = "redis.client"
	// cfg.analyticsRate = globalconfig.AnalyticsRate()
	if internal.BoolEnv("POBP_TRACE_REDIS_ANALYTICS_ENABLED", false) {
		cfg.analyticsRate = 1.0
	} else {
		cfg.analyticsRate = math.NaN()
	}
}

// WithSkipRawCommand reports whether to skip setting the "redis.raw_command" tag
// set up to obfuscate this value and it could contain sensitive information.
func WithSkipRawCommand(skip bool) ClientOption {
	return func(cfg *clientConfig) {
		cfg.skipRaw = skip
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
