package interceptors

import (
	"context"
	"github.com/JECSand/identity-service/pkg/authentication"
	"github.com/JECSand/identity-service/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// InterceptorManager manages the interceptors available to service gRPC requests
type InterceptorManager interface {
	Logger(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error)
	ClientRequestLoggerInterceptor() func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error
	AuthUnaryInterceptor(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error)
	AuthStreamInterceptor(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error)
}

// InterceptorManager struct
type interceptorManager struct {
	logger logging.Logger
	auth   authentication.Authenticator
}

// NewInterceptorManager InterceptorManager constructor
func NewInterceptorManager(logger logging.Logger, auth authentication.Authenticator) *interceptorManager {
	return &interceptorManager{
		logger: logger,
		auth:   auth,
	}
}

// Logger Interceptor
func (im *interceptorManager) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.GrpcMiddlewareAccessLogger(info.FullMethod, time.Since(start), md, err)
	return reply, err
}

// ClientRequestLoggerInterceptor gRPC client interceptor
func (im *interceptorManager) ClientRequestLoggerInterceptor() func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		md, _ := metadata.FromIncomingContext(ctx)
		im.logger.GrpcClientInterceptorLogger(method, req, reply, time.Since(start), md, err)
		return err
	}
}

// AuthUnaryInterceptor intercepts unary gRPC requests and verifies authentication.Session validity
func (im *interceptorManager) AuthUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	im.logger.Info("--> unary auth interceptor: ", info.FullMethod)
	_, err = im.auth.AuthorizeGRPC(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

// AuthStreamInterceptor intercepts unary gRPC requests and verifies authentication.Session validity
func (im *interceptorManager) AuthStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	im.logger.Info("--> stream auth interceptor: ", info.FullMethod)
	_, err = im.auth.AuthorizeGRPC(stream.Context(), info.FullMethod)
	if err != nil {
		return err
	}
	return handler(srv, stream)
}
