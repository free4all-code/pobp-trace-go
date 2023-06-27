
package appsec_test

import (
	"encoding/json"
	"io"
	"net/http"

	"git.proto.group/protoobp/pobp-trace-go/appsec"
	echotrace "git.proto.group/protoobp/pobp-trace-go/contrib/labstack/echo.v4"
	httptrace "git.proto.group/protoobp/pobp-trace-go/contrib/net/http"

	"github.com/labstack/echo/v4"
)

type parsedBodyType struct {
	Value string `json:"value"`
}

func customBodyParser(body io.ReadCloser) (*parsedBodyType, error) {
	var parsedBody parsedBodyType
	err := json.NewDecoder(body).Decode(&parsedBody)
	return &parsedBody, err
}

// Monitor HTTP request parsed body
func ExampleMonitorParsedHTTPBody() {
	mux := httptrace.NewServeMux()
	mux.HandleFunc("/body", func(w http.ResponseWriter, r *http.Request) {
		// Use the SDK to monitor the request's parsed body
		body, err := customBodyParser(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		appsec.MonitorParsedHTTPBody(r.Context(), body)
		w.Write([]byte("Body monitored using AppSec SDK\n"))
	})
	http.ListenAndServe(":8080", mux)
}

// Monitor HTTP request parsed body with a framework customized context type
func ExampleMonitorParsedHTTPBody_CustomContext() {
	r := echo.New()
	r.Use(echotrace.Middleware())
	r.POST("/body", func(c echo.Context) (e error) {
		req := c.Request()
		body, err := customBodyParser(req.Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		// Use the SDK to monitor the request's parsed body
		appsec.MonitorParsedHTTPBody(c.Request().Context(), body)
		return c.String(http.StatusOK, "Body monitored using AppSec SDK")
	})

	r.Start(":8080")
}