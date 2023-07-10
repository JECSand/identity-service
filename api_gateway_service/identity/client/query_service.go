package client

import (
	"context"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/pkg/interceptors"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

const (
	backoffLinear  = 100 * time.Millisecond
	backoffRetries = 3
)

// NewQueryServiceClient constructs and return a new gRPC client connection to the query service
func NewQueryServiceClient(ctx context.Context, cfg *config.Config, im interceptors.InterceptorManager) (*grpc.ClientConn, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffLinear)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
		grpc_retry.WithMax(backoffRetries),
	}
	queryClient, err := grpc.DialContext(
		ctx,
		cfg.Grpc.QueryServicePort,
		grpc.WithUnaryInterceptor(im.ClientRequestLoggerInterceptor()),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.DialContext")
	}
	return queryClient, nil
}
