

package opentracer

import (
	"context"
	"testing"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/internal"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)
	ot := New()
	dd, ok := internal.GetGlobalTracer().(pobptrace.Tracer)
	assert.True(ok)
	ott, ok := ot.(*opentracer)
	assert.True(ok)
	assert.Equal(ott.Tracer, dd)
}

func TestSpanWithContext(t *testing.T) {
	assert := assert.New(t)
	ot, ok := New().(*opentracer)
	assert.True(ok)
	opentracing.SetGlobalTracer(ot)
	want, ctx := opentracing.StartSpanFromContext(context.Background(), "test.operation")
	got, ok := tracer.SpanFromContext(ctx)
	assert.True(ok)
	assert.Equal(got, want.(*span).Span)
}

func TestInjectError(t *testing.T) {
	ot := New()

	for name, tt := range map[string]struct {
		spanContext opentracing.SpanContext
		format      interface{}
		carrier     interface{}
		want        error
	}{
		"ErrInvalidSpanContext": {
			spanContext: internal.NoopSpanContext{},
			format:      opentracing.TextMap,
			carrier:     opentracing.TextMapCarrier(map[string]string{}),
			want:        opentracing.ErrInvalidSpanContext,
		},
		"ErrInvalidCarrier": {
			spanContext: ot.StartSpan("test.operation").Context(),
			format:      opentracing.TextMap,
			carrier:     "invalid-carrier",
			want:        opentracing.ErrInvalidCarrier,
		},
		"ErrUnsupportedFormat": {
			format: "unsupported-format",
			want:   opentracing.ErrUnsupportedFormat,
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := ot.Inject(tt.spanContext, tt.format, tt.carrier)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractError(t *testing.T) {
	ot := New()

	for name, tt := range map[string]struct {
		format  interface{}
		carrier interface{}
		want    error
	}{
		"ErrSpanContextNotFound": {
			format:  opentracing.TextMap,
			carrier: opentracing.TextMapCarrier(nil),
			want:    opentracing.ErrSpanContextNotFound,
		},
		"ErrInvalidCarrier": {
			format:  opentracing.TextMap,
			carrier: "invalid-carrier",
			want:    opentracing.ErrInvalidCarrier,
		},
		"ErrSpanContextCorrupted": {
			format: opentracing.TextMap,
			carrier: opentracing.TextMapCarrier(
				map[string]string{
					tracer.DefaultTraceIDHeader:  "-1",
					tracer.DefaultParentIDHeader: "-1",
					tracer.DefaultPriorityHeader: "not-a-number",
				},
			),
			want: opentracing.ErrSpanContextCorrupted,
		},
		"ErrUnsupportedFormat": {
			format: "unsupported-format",
			want:   opentracing.ErrUnsupportedFormat,
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, got := ot.Extract(tt.format, tt.carrier)
			assert.Equal(t, tt.want, got)
		})
	}
}
