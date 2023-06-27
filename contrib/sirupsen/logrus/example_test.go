
package logrus

import (
	"context"
	"os"
	"time"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	"github.com/sirupsen/logrus"
)

func ExampleHook() {
	// Ensure your tracer is started and stopped
	// Setup logrus, do this once at the beginning of your program
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(&DDContextLogHook{})
	logrus.SetOutput(os.Stdout)

	span, sctx := tracer.StartSpanFromContext(context.Background(), "mySpan")
	defer span.Finish()

	// Pass the current span context to the logger (Time is set for consistency in output here)
	cLog := logrus.WithContext(sctx).WithTime(time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC))
	// Log as desired using the context-aware logger
	cLog.Info("Completed some work!")
	// Output:
	// {"pobp.span_id":0,"pobp.trace_id":0,"level":"info","msg":"Completed some work!","time":"2000-01-01T01:01:01Z"}
}
