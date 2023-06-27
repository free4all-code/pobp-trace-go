

package web_test

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	webtrace "git.proto.group/protoobp/pobp-trace-go/contrib/zenazn/goji.v1/web"
)

func ExampleMiddleware() {
	// Using the Router middleware lets the tracer determine routes for
	// use in a trace's resource name ("GET /user/:id")
	// Otherwise the resource is only the method ("GET", "POST", etc.)
	goji.Use(goji.DefaultMux.Router)
	goji.Use(webtrace.Middleware())
	goji.Get("/hello", func(c web.C, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Why hello there!")
	})
	goji.Serve()
}
