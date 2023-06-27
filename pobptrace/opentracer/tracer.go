

// Package opentracer is in "Maintenance" mode and limited support is offered. Please consider
// using OpenTelemetry or pobptrace/tracer directly. For additional details, please see our Support
// Policy: https://github.com/DataDog/dd-trace-go#support-policy
//
// Package opentracer provides a wrapper on top of the Datadog tracer that can be used with Opentracing.
// It also provides a set of opentracing.StartSpanOption that are specific to Datadog's APM product.
// To use it, simply call "New".
//
// Note that there are currently some small incompatibilities between the Opentracing spec and the Datadog
// APM product, which we are in the process of addressing on the long term. When using Datadog, the
// Opentracing operation name is what is called resource in Datadog's terms and the Opentracing "component"
// tag is Datadog's operation name. Meaning that in order to define (in Opentracing terms) a span that
// has the operation name "/user/profile" and the component "http.request", one would do:
//  opentracing.StartSpan("http.request", opentracer.ResourceName("/user/profile"))
//
// Some libraries and frameworks are supported out-of-the-box by using our integrations. You can see a list
// of supported integrations here: https://godoc.org/git.proto.group/protoobp/pobp-trace-go/contrib. They are fully
// compatible with the Opentracing implementation.
package opentracer

import (
	"context"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/internal"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	opentracing "github.com/opentracing/opentracing-go"
)

// New creates, instantiates and returns an Opentracing compatible version of the
// Datadog tracer using the provided set of options.
func New(opts ...tracer.StartOption) opentracing.Tracer {
	tracer.Start(opts...)
	return &opentracer{internal.GetGlobalTracer()}
}

var _ opentracing.Tracer = (*opentracer)(nil)

// opentracer implements opentracing.Tracer on top of pobptrace.Tracer.
type opentracer struct{ pobptrace.Tracer }

// StartSpan implements opentracing.Tracer.
func (t *opentracer) StartSpan(operationName string, options ...opentracing.StartSpanOption) opentracing.Span {
	var sso opentracing.StartSpanOptions
	for _, o := range options {
		o.Apply(&sso)
	}
	opts := []pobptrace.StartSpanOption{tracer.StartTime(sso.StartTime)}
	for _, ref := range sso.References {
		if v, ok := ref.ReferencedContext.(pobptrace.SpanContext); ok {
			// opentracing.ChildOfRef and opentracing.FollowsFromRef will both be represented as
			// children because Datadog APM does not have a concept of FollowsFrom references.
			opts = append(opts, tracer.ChildOf(v))
			break // can only have one parent
		}
	}
	for k, v := range sso.Tags {
		opts = append(opts, tracer.Tag(k, v))
	}
	return &span{
		Span:       t.Tracer.StartSpan(operationName, opts...),
		opentracer: t,
	}
}

// Inject implements opentracing.Tracer.
func (t *opentracer) Inject(ctx opentracing.SpanContext, format interface{}, carrier interface{}) error {
	sctx, ok := ctx.(pobptrace.SpanContext)
	if !ok {
		return opentracing.ErrUnsupportedFormat
	}
	switch format {
	case opentracing.TextMap, opentracing.HTTPHeaders:
		return translateError(t.Tracer.Inject(sctx, carrier))
	default:
		return opentracing.ErrUnsupportedFormat
	}
}

// Extract implements opentracing.Tracer.
func (t *opentracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	switch format {
	case opentracing.TextMap, opentracing.HTTPHeaders:
		sctx, err := t.Tracer.Extract(carrier)
		return sctx, translateError(err)
	default:
		return nil, opentracing.ErrUnsupportedFormat
	}
}

var _ opentracing.TracerContextWithSpanExtension = (*opentracer)(nil)

// ContextWithSpan implements opentracing.TracerContextWithSpanExtension.
func (t *opentracer) ContextWithSpanHook(ctx context.Context, openSpan opentracing.Span) context.Context {
	ddSpan, ok := openSpan.(*span)
	if !ok {
		return ctx
	}
	return tracer.ContextWithSpan(ctx, ddSpan.Span)
}

func translateError(err error) error {
	switch err {
	case tracer.ErrSpanContextNotFound:
		return opentracing.ErrSpanContextNotFound
	case tracer.ErrInvalidCarrier:
		return opentracing.ErrInvalidCarrier
	case tracer.ErrInvalidSpanContext:
		return opentracing.ErrInvalidSpanContext
	case tracer.ErrSpanContextCorrupted:
		return opentracing.ErrSpanContextCorrupted
	default:
		return err
	}
}
