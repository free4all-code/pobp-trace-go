

package http_test

import (
	"net/http"

	httptrace "git.proto.group/protoobp/pobp-trace-go/contrib/net/http"
)

func Example() {
	mux := httptrace.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!\n"))
	})
	http.ListenAndServe(":8080", mux)
}

func Example_withServiceName() {
	mux := httptrace.NewServeMux(httptrace.WithServiceName("my-service"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!\n"))
	})
	http.ListenAndServe(":8080", mux)
}

func ExampleTraceAndServe() {
	mux := http.NewServeMux()
	mux.Handle("/", traceMiddleware(mux, http.HandlerFunc(Index)))
	http.ListenAndServe(":8080", mux)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}

func traceMiddleware(mux *http.ServeMux, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, route := mux.Handler(r)
		resource := r.Method + " " + route
		httptrace.TraceAndServe(next, w, r, &httptrace.ServeConfig{
			Service:     "http.router",
			Resource:    resource,
			QueryParams: true,
		})
	})
}
