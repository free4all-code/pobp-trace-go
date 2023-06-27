
// Package logrus provides a log/span correlation hook for the sirupsen/logrus package (https://github.com/sirupsen/logrus).
package logrus

import (
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	"github.com/sirupsen/logrus"
)

// DDContextLogHook ensures that any span in the log context is correlated to log output.
type DDContextLogHook struct{}

// Levels implements logrus.Hook interface, this hook applies to all defined levels
func (d *DDContextLogHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel}
}

// Fire implements logrus.Hook interface, attaches trace and span details found in entry context
func (d *DDContextLogHook) Fire(e *logrus.Entry) error {
	span, found := tracer.SpanFromContext(e.Context)
	if !found {
		return nil
	}
	e.Data["pobp.trace_id"] = span.Context().TraceID()
	e.Data["pobp.span_id"] = span.Context().SpanID()
	return nil
}
