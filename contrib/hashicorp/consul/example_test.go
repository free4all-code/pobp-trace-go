

package consul

import (
	"context"
	"fmt"
	"log"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	consul "github.com/hashicorp/consul/api"
)

// Here's an example illustrating a simple use case for interacting with consul with tracing enabled.
func Example() {
	// Get a new Consul client
	client, err := NewClient(consul.DefaultConfig(), WithServiceName("consul.example"))
	if err != nil {
		log.Fatal(err)
	}

	// Optionally, create a new root span
	root, ctx := tracer.StartSpanFromContext(context.Background(), "root_span",
		tracer.SpanType(ext.SpanTypeConsul),
		tracer.ServiceName("example"),
	)
	defer root.Finish()
	client = client.WithContext(ctx)

	// Get a handle to the KV API
	kv := client.KV()

	// PUT a new KV pair
	p := &consul.KVPair{Key: "test", Value: []byte("1000")}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Lookup the pair
	pair, _, err := kv.Get("test", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v: %s\n", pair.Key, pair.Value)
	// Output:
	// test: 1000
}
