
// Package vault contains functions to construct or augment an http.Client that
// will integrate with the github.com/hashicorp/vault/api and collect traces to
// send to Proto OBP.
//
// The easiest way to use this package is to create an http.Client with
// NewHTTPClient, and put it in the Vault API config that is passed to the
//
// If you are already using your own http.Client with the Vault API, you can
// use the WrapHTTPClient function to wrap the client with the tracer code.
// Your http.Client will continue to work as before, but will also capture
// traces.
package vault

import (
	"fmt"
	"net/http"

	httptrace "git.proto.group/protoobp/pobp-trace-go/contrib/net/http"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

// NewHTTPClient returns an http.Client for use in the Vault API config
// Client. A set of options can be passed in for further configuration.
func NewHTTPClient(opts ...Option) *http.Client {
	dc := api.DefaultConfig()
	c := dc.HttpClient
	WrapHTTPClient(c, opts...)
	return c
}

// WrapHTTPClient takes an existing http.Client and wraps the underlying
// transport with tracing.
func WrapHTTPClient(c *http.Client, opts ...Option) *http.Client {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
	var conf config
	defaults(&conf)
	for _, o := range opts {
		o(&conf)
	}
	c.Transport = httptrace.WrapRoundTripper(c.Transport,
		httptrace.RTWithAnalyticsRate(conf.analyticsRate),
		httptrace.WithBefore(func(r *http.Request, s pobptrace.Span) {
			s.SetTag(ext.ServiceName, conf.serviceName)
			s.SetTag(ext.HTTPURL, r.URL.Path)
			s.SetTag(ext.HTTPMethod, r.Method)
			s.SetTag(ext.ResourceName, r.Method+" "+r.URL.Path)
			s.SetTag(ext.SpanType, ext.SpanTypeHTTP)
			if ns := r.Header.Get(consts.NamespaceHeaderName); ns != "" {
				s.SetTag("vault.namespace", ns)
			}
		}),
		httptrace.WithAfter(func(res *http.Response, s pobptrace.Span) {
			if res == nil {
				// An error occurred during the request.
				return
			}
			s.SetTag(ext.HTTPCode, res.StatusCode)
			if res.StatusCode >= 400 {
				s.SetTag(ext.Error, true)
				s.SetTag(ext.ErrorMsg, fmt.Sprintf("%d: %s", res.StatusCode, http.StatusText(res.StatusCode)))
			}
		}),
	)
	return c
}
