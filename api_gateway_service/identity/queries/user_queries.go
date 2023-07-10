package queries

import (
	"context"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	queryService "github.com/JECSand/identity-service/query_service/protos/user_query"
	"github.com/opentracing/opentracing-go"
)

type GetUserByIdHandler interface {
	Handle(ctx context.Context, query *GetUserByIdQuery) (*dto.UserResponse, error)
}

type getUserByIdHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient queryService.QueryServiceClient
}

func NewGetUserByIdHandler(log logging.Logger, cfg *config.Config, rsClient queryService.QueryServiceClient) *getUserByIdHandler {
	return &getUserByIdHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (q *getUserByIdHandler) Handle(ctx context.Context, query *GetUserByIdQuery) (*dto.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getUserByIdHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := q.rsClient.GetUserById(ctx, &queryService.GetUserByIdReq{ID: query.ID.String()})
	if err != nil {
		return nil, err
	}
	return dto.UserResponseFromGrpc(res.GetUser()), nil
}

// SearchUserHandler ...
type SearchUserHandler interface {
	Handle(ctx context.Context, query *SearchUserQuery) (*dto.UsersListResponse, error)
}

type searchUserHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient queryService.QueryServiceClient
}

func NewSearchUserHandler(log logging.Logger, cfg *config.Config, rsClient queryService.QueryServiceClient) *searchUserHandler {
	return &searchUserHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (s *searchUserHandler) Handle(ctx context.Context, query *SearchUserQuery) (*dto.UsersListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "searchUserHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := s.rsClient.SearchUser(ctx, &queryService.SearchReq{
		Search: query.Text,
		Page:   int64(query.Pagination.GetPage()),
		Size:   int64(query.Pagination.GetSize()),
	})
	if err != nil {
		return nil, err
	}
	return dto.UsersListResponseFromGrpc(res), nil
}
