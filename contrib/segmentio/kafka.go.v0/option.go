

package kafka

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

// An Option customizes the config.
type Option func(cfg *config)

func newConfig(opts ...Option) *config {
	cfg := &config{
		consumerServiceName: "kafka",
		producerServiceName: "kafka",
		// analyticsRate: globalconfig.AnalyticsRate(),
		analyticsRate: math.NaN(),
	}
	if internal.BoolEnv("POBP_TRACE_KAFKA_ANALYTICS_ENABLED", false) {
		cfg.analyticsRate = 1.0
	}
	if svc := globalconfig.ServiceName(); svc != "" {
		cfg.consumerServiceName = svc
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithServiceName sets the config service name to serviceName.
func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.consumerServiceName = serviceName
		cfg.producerServiceName = serviceName
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
