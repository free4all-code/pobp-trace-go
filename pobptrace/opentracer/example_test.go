

package opentracer_test

import (
	opentracing "github.com/opentracing/opentracing-go"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/opentracer"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
)

func Example() {
	// Start a Datadog tracer, optionally providing a set of options,
	// returning an opentracing.Tracer which wraps it.
	t := opentracer.New(tracer.WithAgentAddr("host:port"))

	// Use it with the Opentracing API. The (already started) Datadog tracer
	// may be used in parallel with the Opentracing API if desired.
	opentracing.SetGlobalTracer(t)
}
