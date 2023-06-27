

package grpc

import (
	"encoding/json"
	"net"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation/grpcsec"
	"git.proto.group/protoobp/pobp-trace-go/internal/appsec/dyngo/instrumentation/httpsec"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// UnaryHandler wrapper to use when AppSec is enabled to monitor its execution.
func appsecUnaryHandlerMiddleware(span pobptrace.Span, handler grpc.UnaryHandler) grpc.UnaryHandler {
	httpsec.SetAppSecTags(span)
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		op := grpcsec.StartHandlerOperation(grpcsec.HandlerOperationArgs{Metadata: md}, nil)
		defer func() {
			events := op.Finish(grpcsec.HandlerOperationRes{})
			instrumentation.SetTags(span, op.Tags())
			if len(events) == 0 {
				return
			}
			setAppSecTags(ctx, span, events)
		}()
		defer grpcsec.StartReceiveOperation(grpcsec.ReceiveOperationArgs{}, op).Finish(grpcsec.ReceiveOperationRes{Message: req})
		return handler(ctx, req)
	}
}

// StreamHandler wrapper to use when AppSec is enabled to monitor its execution.
func appsecStreamHandlerMiddleware(span pobptrace.Span, handler grpc.StreamHandler) grpc.StreamHandler {
	httpsec.SetAppSecTags(span)
	return func(srv interface{}, stream grpc.ServerStream) error {
		md, _ := metadata.FromIncomingContext(stream.Context())
		op := grpcsec.StartHandlerOperation(grpcsec.HandlerOperationArgs{Metadata: md}, nil)
		defer func() {
			events := op.Finish(grpcsec.HandlerOperationRes{})
			instrumentation.SetTags(span, op.Tags())
			if len(events) == 0 {
				return
			}
			setAppSecTags(stream.Context(), span, events)
		}()
		return handler(srv, appsecServerStream{ServerStream: stream, handlerOperation: op})
	}
}

type appsecServerStream struct {
	grpc.ServerStream
	handlerOperation *grpcsec.HandlerOperation
}

// RecvMsg implements grpc.ServerStream interface method to monitor its
// execution with AppSec.
func (ss appsecServerStream) RecvMsg(m interface{}) error {
	op := grpcsec.StartReceiveOperation(grpcsec.ReceiveOperationArgs{}, ss.handlerOperation)
	defer func() {
		op.Finish(grpcsec.ReceiveOperationRes{Message: m})
	}()
	return ss.ServerStream.RecvMsg(m)
}

// Set the AppSec tags when security events were found.
func setAppSecTags(ctx context.Context, span pobptrace.Span, events []json.RawMessage) {
	md, _ := metadata.FromIncomingContext(ctx)
	var addr net.Addr
	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr
	}
	grpcsec.SetSecurityEventTags(span, events, addr, md)
}
