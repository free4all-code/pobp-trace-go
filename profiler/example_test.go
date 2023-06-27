

package profiler_test

import (
	"log"

	"git.proto.group/protoobp/pobp-trace-go/profiler"
)

// This example illustrates how to run (and later stop) the Datadog Profiler.
func Example() {
	err := profiler.Start(
		profiler.WithService("users-db"),
		profiler.WithEnv("staging"),
		profiler.WithTags("version:1.2.0"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()

	// ...
}
