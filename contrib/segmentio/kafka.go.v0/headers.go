

package kafka

import (
	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	"github.com/segmentio/kafka-go"
)

// A messageCarrier implements TextMapReader/TextMapWriter for extracting/injecting traces on a kafka.Message
type messageCarrier struct {
	msg *kafka.Message
}

var _ interface {
	tracer.TextMapReader
	tracer.TextMapWriter
} = (*messageCarrier)(nil)

// ForeachKey conforms to the TextMapReader interface.
func (c messageCarrier) ForeachKey(handler func(key, val string) error) error {
	for _, h := range c.msg.Headers {
		err := handler(h.Key, string(h.Value))
		if err != nil {
			return err
		}
	}
	return nil
}

// Set implements TextMapWriter
func (c messageCarrier) Set(key, val string) {
	// ensure uniqueness of keys
	for i := 0; i < len(c.msg.Headers); i++ {
		if string(c.msg.Headers[i].Key) == key {
			c.msg.Headers = append(c.msg.Headers[:i], c.msg.Headers[i+1:]...)
			i--
		}
	}
	c.msg.Headers = append(c.msg.Headers, kafka.Header{
		Key:   key,
		Value: []byte(val),
	})
}

// ExtractSpanContext retrieves the SpanContext from a kafka.Message
func ExtractSpanContext(msg kafka.Message) (pobptrace.SpanContext, error) {
	return tracer.Extract(messageCarrier{&msg})
}
