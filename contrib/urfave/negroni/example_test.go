
package negroni_test

import (
	"net/http"

	"github.com/urfave/negroni"

	negronitrace "git.proto.group/protoobp/pobp-trace-go/contrib/urfave/negroni"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello World!\n"))
}

func Example() {
	// Start the tracer
	tracer.Start()
	defer tracer.Stop()

	// Create a negroni Router
	n := negroni.New()

	// Use the tracer middleware with the default service name "negroni.router".
	n.Use(negronitrace.Middleware())

	// Set up some endpoints.
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	n.UseHandler(mux)

	// And start gathering request traces
	http.ListenAndServe(":8080", n)
}

func Example_withServiceName() {
	// Start the tracer
	tracer.Start()
	defer tracer.Stop()

	// Create a negroni Router
	n := negroni.New()

	// Use the tracer middleware with your desired service name.
	n.Use(negronitrace.Middleware(negronitrace.WithServiceName("negroni-server")))

	// Set up some endpoints.
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	n.UseHandler(mux)

	// And start gathering request traces
	http.ListenAndServe(":8080", n)
}
