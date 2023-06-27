
package logrus

import (
	"context"
	"testing"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFire(t *testing.T) {
	tracer.Start()
	defer tracer.Stop()
	_, sctx := tracer.StartSpanFromContext(context.Background(), "testSpan", tracer.WithSpanID(1234))

	hook := &DDContextLogHook{}
	e := logrus.NewEntry(logrus.New())
	e.Context = sctx
	err := hook.Fire(e)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1234), e.Data["pobp.trace_id"])
	assert.Equal(t, uint64(1234), e.Data["pobp.span_id"])
}
