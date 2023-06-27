

package pg

import (
	"context"
	"math"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/log"

	"github.com/go-pg/pg/v10"
)

// Wrap augments the given DB with tracing.
func Wrap(db *pg.DB, opts ...Option) {
	cfg := new(config)
	defaults(cfg)
	for _, opt := range opts {
		opt(cfg)
	}
	log.Debug("contrib/go-pg/pg.v10: Wrapping Database")
	db.AddQueryHook(&queryHook{cfg: cfg})
}

type queryHook struct {
	cfg *config
}

// BeforeQuery implements pg.QueryHook.
func (h *queryHook) BeforeQuery(ctx context.Context, qe *pg.QueryEvent) (context.Context, error) {
	query, err := qe.UnformattedQuery()
	if err != nil {
		query = []byte("unknown")
	}

	opts := []pobptrace.StartSpanOption{
		tracer.SpanType(ext.SpanTypeSQL),
		tracer.ResourceName(string(query)),
		tracer.ServiceName(h.cfg.serviceName),
	}
	if !math.IsNaN(h.cfg.analyticsRate) {
		opts = append(opts, tracer.Tag(ext.EventSampleRate, h.cfg.analyticsRate))
	}
	_, ctx = tracer.StartSpanFromContext(ctx, "go-pg", opts...)
	return ctx, qe.Err
}

// AfterQuery implements pg.QueryHook
func (h *queryHook) AfterQuery(ctx context.Context, qe *pg.QueryEvent) error {
	if span, ok := tracer.SpanFromContext(ctx); ok {
		span.Finish(tracer.WithError(qe.Err))
	}

	return qe.Err
}
