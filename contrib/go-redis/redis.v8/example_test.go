

package redis_test

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	redistrace "git.proto.group/protoobp/pobp-trace-go/contrib/go-redis/redis.v8"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
)

// To start tracing Redis, simply create a new client using the library and continue
// using as you normally would.
func Example() {
	ctx := context.Background()
	// create a new Client
	opts := &redis.Options{Addr: "127.0.0.1", Password: "", DB: 0}
	c := redistrace.NewClient(opts)

	// any action emits a span
	c.Set(ctx, "test_key", "test_value", 0)

	// optionally, create a new root span
	root, ctx := tracer.StartSpanFromContext(context.Background(), "parent.request",
		tracer.SpanType(ext.SpanTypeRedis),
		tracer.ServiceName("web"),
		tracer.ResourceName("/home"),
	)

	// commit further commands, which will inherit from the parent in the context.
	c.Set(ctx, "food", "cheese", 0)
	root.Finish()
}

// You can also trace Redis Pipelines. Simply use as usual and the traces will be
// automatically picked up by the underlying implementation.
func Example_pipeliner() {
	ctx := context.Background()
	// create a client
	opts := &redis.Options{Addr: "127.0.0.1", Password: "", DB: 0}
	c := redistrace.NewClient(opts, redistrace.WithServiceName("my-redis-service"))

	// open the pipeline
	pipe := c.Pipeline()

	// submit some commands
	pipe.Incr(ctx, "pipeline_counter")
	pipe.Expire(ctx, "pipeline_counter", time.Hour)

	// execute with trace
	pipe.Exec(ctx)
}
