

package tracer

import (
	"context"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/internal"
)

type contextKey struct{}

var activeSpanKey = contextKey{}

// ContextWithSpan returns a copy of the given context which includes the span s.
func ContextWithSpan(ctx context.Context, s Span) context.Context {
	return context.WithValue(ctx, activeSpanKey, s)
}

// SpanFromContext returns the span contained in the given context. A second return
// value indicates if a span was found in the context. If no span is found, a no-op
// span is returned.
func SpanFromContext(ctx context.Context) (Span, bool) {
	if ctx == nil {
		return &internal.NoopSpan{}, false
	}
	v := ctx.Value(activeSpanKey)
	if s, ok := v.(pobptrace.Span); ok {
		return s, true
	}
	return &internal.NoopSpan{}, false
}

// StartSpanFromContext returns a new span with the given operation name and options. If a span
// is found in the context, it will be used as the parent of the resulting span. If the ChildOf
// option is passed, the span from context will take precedence over it as the parent span.
func StartSpanFromContext(ctx context.Context, operationName string, opts ...StartSpanOption) (Span, context.Context) {
	// copy opts in case the caller reuses the slice in parallel
	// we will add at least 1, at most 2 items
	optsLocal := make([]StartSpanOption, len(opts), len(opts)+2)
	copy(optsLocal, opts)

	if ctx == nil {
		// default to context.Background() to avoid panics on Go >= 1.15
		ctx = context.Background()
	} else if s, ok := SpanFromContext(ctx); ok {
		optsLocal = append(optsLocal, ChildOf(s.Context()))
	}
	optsLocal = append(optsLocal, withContext(ctx))
	s := StartSpan(operationName, optsLocal...)
	if span, ok := s.(*span); ok && span.pprofCtxActive != nil {
		// If pprof labels were applied for this span, use the derived ctx that
		// includes them. Otherwise a child of this span wouldn't be able to
		// correctly restore the labels of its parent when it finishes.
		ctx = span.pprofCtxActive
	}
	return s, ContextWithSpan(ctx, s)
}
