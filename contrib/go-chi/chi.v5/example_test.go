

package chi_test

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	chitrace "git.proto.group/protoobp/pobp-trace-go/contrib/go-chi/chi.v5"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}

func Example() {
	// Start the tracer
	tracer.Start()
	defer tracer.Stop()

	// Create a chi Router
	router := chi.NewRouter()

	// Use the tracer middleware with the default service name "chi.router".
	router.Use(chitrace.Middleware())

	// Set up some endpoints.
	router.Get("/", handler)

	// And start gathering request traces
	http.ListenAndServe(":8080", router)
}

func Example_withServiceName() {
	// Start the tracer
	tracer.Start()
	defer tracer.Stop()

	// Create a chi Router
	router := chi.NewRouter()

	// Use the tracer middleware with your desired service name.
	router.Use(chitrace.Middleware(chitrace.WithServiceName("chi-server")))

	// Set up some endpoints.
	router.Get("/", handler)

	// And start gathering request traces
	http.ListenAndServe(":8080", router)
}
