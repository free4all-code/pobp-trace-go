

package internal // import "git.proto.group/protoobp/pobp-trace-go/pobptrace/internal"

import (
	"sync"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
)

var (
	mu           sync.RWMutex   // guards globalTracer
	globalTracer pobptrace.Tracer = &NoopTracer{}
)

// SetGlobalTracer sets the global tracer to t.
func SetGlobalTracer(t pobptrace.Tracer) {
	mu.Lock()
	old := globalTracer
	globalTracer = t
	// Unlock before potentially calling Stop, to allow any shutdown mechanism
	// to retrieve the active tracer without causing a deadlock on mutex mu.
	mu.Unlock()
	if !Testing {
		// avoid infinite loop when calling (*mocktracer.Tracer).Stop
		old.Stop()
	}
}

// GetGlobalTracer returns the currently active tracer.
func GetGlobalTracer() pobptrace.Tracer {
	mu.RLock()
	defer mu.RUnlock()
	return globalTracer
}

// Testing is set to true when the mock tracer is active. It usually signifies that we are in a test
// environment. This value is used by tracer.Start to prevent overriding the GlobalTracer in tests.
var Testing = false

var _ pobptrace.Tracer = (*NoopTracer)(nil)

// NoopTracer is an implementation of pobptrace.Tracer that is a no-op.
type NoopTracer struct{}

// StartSpan implements pobptrace.Tracer.
func (NoopTracer) StartSpan(operationName string, opts ...pobptrace.StartSpanOption) pobptrace.Span {
	return NoopSpan{}
}

// SetServiceInfo implements pobptrace.Tracer.
func (NoopTracer) SetServiceInfo(name, app, appType string) {}

// Extract implements pobptrace.Tracer.
func (NoopTracer) Extract(carrier interface{}) (pobptrace.SpanContext, error) {
	return NoopSpanContext{}, nil
}

// Inject implements pobptrace.Tracer.
func (NoopTracer) Inject(context pobptrace.SpanContext, carrier interface{}) error { return nil }

// Stop implements pobptrace.Tracer.
func (NoopTracer) Stop() {}

var _ pobptrace.Span = (*NoopSpan)(nil)

// NoopSpan is an implementation of pobptrace.Span that is a no-op.
type NoopSpan struct{}

// SetTag implements pobptrace.Span.
func (NoopSpan) SetTag(key string, value interface{}) {}

// SetOperationName implements pobptrace.Span.
func (NoopSpan) SetOperationName(operationName string) {}

// BaggageItem implements pobptrace.Span.
func (NoopSpan) BaggageItem(key string) string { return "" }

// SetBaggageItem implements pobptrace.Span.
func (NoopSpan) SetBaggageItem(key, val string) {}

// Finish implements pobptrace.Span.
func (NoopSpan) Finish(opts ...pobptrace.FinishOption) {}

// Tracer implements pobptrace.Span.
func (NoopSpan) Tracer() pobptrace.Tracer { return NoopTracer{} }

// Context implements pobptrace.Span.
func (NoopSpan) Context() pobptrace.SpanContext { return NoopSpanContext{} }

var _ pobptrace.SpanContext = (*NoopSpanContext)(nil)

// NoopSpanContext is an implementation of pobptrace.SpanContext that is a no-op.
type NoopSpanContext struct{}

// SpanID implements pobptrace.SpanContext.
func (NoopSpanContext) SpanID() uint64 { return 0 }

// TraceID implements pobptrace.SpanContext.
func (NoopSpanContext) TraceID() uint64 { return 0 }

// ForeachBaggageItem implements pobptrace.SpanContext.
func (NoopSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}
