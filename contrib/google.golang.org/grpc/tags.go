

package grpc

// Tags used for gRPC
const (
	tagMethodName     = "grpc.method.name"
	tagMethodKind     = "grpc.method.kind"
	tagCode           = "grpc.code"
	tagMetadataPrefix = "grpc.metadata."
	tagRequest        = "grpc.request"
)

const (
	methodKindUnary        = "unary"
	methodKindClientStream = "client_streaming"
	methodKindServerStream = "server_streaming"
	methodKindBidiStream   = "bidi_streaming"
)
