

// Package fiber provides tracing functions for tracing the fiber package (https://github.com/gofiber/fiber).
package fiber // import "git.proto.group/protoobp/pobp-trace-go/contrib/gofiber/fiber.v2"

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/log"

	"github.com/gofiber/fiber/v2"
)

// Middleware returns middleware that will trace incoming requests.
func Middleware(opts ...Option) func(c *fiber.Ctx) error {
	cfg := new(config)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	log.Debug("gofiber/fiber.v2: Middleware: %#v", cfg)
	return func(c *fiber.Ctx) error {
		opts := []pobptrace.StartSpanOption{
			tracer.SpanType(ext.SpanTypeWeb),
			tracer.ServiceName(cfg.serviceName),
			tracer.Tag(ext.HTTPMethod, c.Method()),
			tracer.Tag(ext.HTTPURL, string(c.Request().URI().PathOriginal())),
			tracer.Measured(),
		}
		if !math.IsNaN(cfg.analyticsRate) {
			opts = append(opts, tracer.Tag(ext.EventSampleRate, cfg.analyticsRate))
		}

		opts = append(opts, cfg.spanOpts...)
		span, ctx := tracer.StartSpanFromContext(c.Context(), "http.request", opts...)

		defer span.Finish()

		// pass the span through the request UserContext
		c.SetUserContext(ctx)

		// pass the execution down the line
		err := c.Next()

		span.SetTag(ext.ResourceName, cfg.resourceNamer(c))

		status := c.Response().StatusCode()
		// on the off chance we don't yet have a status after the rest of the things have run
		if status == 0 {
			// 0 - means we do not have a status code at this point
			// in case the response was returned by a middleware without one
			status = http.StatusOK
		}
		span.SetTag(ext.HTTPCode, strconv.Itoa(status))

		if err != nil {
			span.SetTag(ext.Error, err)
		} else if cfg.isStatusError(status) {
			// mark 5xx server error
			span.SetTag(ext.Error, fmt.Errorf("%d: %s", status, http.StatusText(status)))
		}
		return err
	}
}