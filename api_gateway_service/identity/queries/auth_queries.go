package queries

import (
	"context"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	authQueryService "github.com/JECSand/identity-service/query_service/protos/auth_query"
	"github.com/opentracing/opentracing-go"
)

type AuthenticateHandler interface {
	Handle(ctx context.Context, query *AuthenticateQuery) (*dto.AuthenticateResponse, error)
}

type authenticateHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient authQueryService.AuthQueryServiceClient
}

func NewAuthenticateHandler(log logging.Logger, cfg *config.Config, rsClient authQueryService.AuthQueryServiceClient) *authenticateHandler {
	return &authenticateHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (q *authenticateHandler) Handle(ctx context.Context, query *AuthenticateQuery) (*dto.AuthenticateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authenticateHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := q.rsClient.Authenticate(ctx, &authQueryService.AuthenticateReq{
		Email:    query.Email,
		Password: query.Password,
	})
	if err != nil {
		return nil, err
	}
	return dto.AuthenticateResponseResponseFromGrpc(res), nil
}

// ValidateHandler ...
type ValidateHandler interface {
	Handle(ctx context.Context, query *ValidateQuery) (*dto.ValidateResponse, error)
}

type validateHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient authQueryService.AuthQueryServiceClient
}

func NewValidateHandler(log logging.Logger, cfg *config.Config, rsClient authQueryService.AuthQueryServiceClient) *validateHandler {
	return &validateHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (s *validateHandler) Handle(ctx context.Context, query *ValidateQuery) (*dto.ValidateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "validateHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := s.rsClient.Validate(ctx, &authQueryService.ValidateReq{
		UserID:         query.UserID,
		AccessToken:    query.AccessToken,
		ValidationType: int64(query.ValidationType.EnumIndex()),
	})
	if err != nil {
		return nil, err
	}
	return dto.ValidateResponseResponseFromGrpc(res), nil
}
