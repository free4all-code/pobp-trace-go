
package elastic_test

import (
	"context"
	"log"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"

	elastictrace "git.proto.group/protoobp/pobp-trace-go/contrib/elastic/go-elasticsearch.v6"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
)

func Example_v7() {
	cfg := elasticsearch.Config{
		Transport: elastictrace.NewRoundTripper(elastictrace.WithServiceName("my-es-service")),
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	_, err = esapi.IndexRequest{
		Index:        "twitter",
		DocumentID:   "1",
		DocumentType: "tweet",
		Body:         strings.NewReader(`{"user": "test", "message": "hello"}`),
	}.Do(context.Background(), es)

	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}

	// Use a context to pass information down the call chain
	root, ctx := tracer.StartSpanFromContext(context.Background(), "parent.request",
		tracer.ServiceName("web"),
		tracer.ResourceName("/tweet/1"),
	)

	_, err = esapi.GetRequest{
		Index:        "twitter",
		DocumentID:   "1",
		DocumentType: "tweet",
	}.Do(ctx, es)

	if err != nil {
		log.Fatalf("Error getting index: %s", err)
	}

	root.Finish()

}
