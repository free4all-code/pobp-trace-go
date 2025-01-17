

package sarama

import (
	"math"

	"git.proto.group/protoobp/pobp-trace-go/internal"
	"git.proto.group/protoobp/pobp-trace-go/internal/globalconfig"
)

type config struct {
	consumerServiceName string
	producerServiceName string
	analyticsRate       float64
}

func defaults(cfg *config) {
	cfg.producerServiceName = "kafka"
	cfg.consumerServiceName = "kafka"
	if svc := globalconfig.ServiceName(); svc != "" {
		cfg.consumerServiceName = svc
	}
	// cfg.analyticsRate = globalconfig.AnalyticsRate()
	if internal.BoolEnv("POBP_TRACE_SARAMA_ANALYTICS_ENABLED", false) {
		cfg.analyticsRate = 1.0
	} else {
		cfg.analyticsRate = math.NaN()
	}
}

// An Option is used to customize the config for the sarama tracer.
type Option func(cfg *config)

// WithServiceName sets the given service name for the intercepted client.
func WithServiceName(name string) Option {
	return func(cfg *config) {
		cfg.consumerServiceName = name
		cfg.producerServiceName = name
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
