

package fiber_test

import (
	fibertrace "git.proto.group/protoobp/pobp-trace-go/contrib/gofiber/fiber.v2"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	"github.com/gofiber/fiber/v2"
)

func Example() {
	// Start the tracer
	tracer.Start()
	defer tracer.Stop()

	// Create a fiber v2 Router
	router := fiber.New()

	// Use the tracer middleware with the default service name "fiber".
	router.Use(fibertrace.Middleware())

	// Set up some endpoints.
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	// And start gathering request traces
	router.Listen(":8080")
}

func Example_withServiceName() {
	// Start the tracer
	tracer.Start()
	defer tracer.Stop()

	// Create a fiber v2 Router
	router := fiber.New()

	// Use the tracer middleware with your desired service name.
	router.Use(fibertrace.Middleware(fibertrace.WithServiceName("fiber")))

	// Set up some endpoints.
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	// And start gathering request traces
	router.Listen(":8080")
}
