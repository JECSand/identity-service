package authentication

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetTokenFromContext parses an auth token from content metadata
func GetTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	return values[0], nil
}

// AttachTokenToContext inputs ctx and an auth token and returns ctx with token attached
func AttachTokenToContext(ctx context.Context, authToken string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", authToken)
}
