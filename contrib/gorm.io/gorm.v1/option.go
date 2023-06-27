

package gorm

import (
	"math"

	"git.proto.group/protoobp/pobp-trace-go/internal"
)

type config struct {
	serviceName   string
	analyticsRate float64
	dsn           string
	errCheck      func(err error) bool
}

// Option represents an option that can be passed to Register, Open or OpenDB.
type Option func(*config)

func defaults(cfg *config) {
	cfg.serviceName = "gorm.db"
	// cfg.analyticsRate = globalconfig.AnalyticsRate()
	if internal.BoolEnv("POBP_TRACE_GORM_ANALYTICS_ENABLED", false) {
		cfg.analyticsRate = 1.0
	} else {
		cfg.analyticsRate = math.NaN()
	}
	cfg.errCheck = func(error) bool { return true }
}

// WithServiceName sets the given service name when registering a driver,
// or opening a database connection.
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

// WithErrorCheck specifies a function fn which determines whether the passed
// error should be marked as an error. The fn is called whenever a gorm operation
// finishes
func WithErrorCheck(fn func(err error) bool) Option {
	return func(cfg *config) {
		cfg.errCheck = fn
	}
}
