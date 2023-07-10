package queries

import (
	"context"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	groupQueryService "github.com/JECSand/identity-service/query_service/protos/group_query"
	"github.com/opentracing/opentracing-go"
)

type GetGroupByIdHandler interface {
	Handle(ctx context.Context, query *GetGroupByIdQuery) (*dto.GroupResponse, error)
}

type getGroupByIdHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient groupQueryService.GroupQueryServiceClient
}

func NewGetGroupByIdHandler(log logging.Logger, cfg *config.Config, rsClient groupQueryService.GroupQueryServiceClient) *getGroupByIdHandler {
	return &getGroupByIdHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (q *getGroupByIdHandler) Handle(ctx context.Context, query *GetGroupByIdQuery) (*dto.GroupResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getGroupByIdHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := q.rsClient.GetGroupById(ctx, &groupQueryService.GetGroupByIdReq{ID: query.ID.String()})
	if err != nil {
		return nil, err
	}
	return dto.GroupResponseFromGrpc(res.GetGroup()), nil
}

// SearchGroupHandler ...
type SearchGroupHandler interface {
	Handle(ctx context.Context, query *SearchGroupQuery) (*dto.GroupsListResponse, error)
}

type searchGroupHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient groupQueryService.GroupQueryServiceClient
}

func NewSearchGroupHandler(log logging.Logger, cfg *config.Config, rsClient groupQueryService.GroupQueryServiceClient) *searchGroupHandler {
	return &searchGroupHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (s *searchGroupHandler) Handle(ctx context.Context, query *SearchGroupQuery) (*dto.GroupsListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "searchGroupHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := s.rsClient.SearchGroup(ctx, &groupQueryService.SearchGroupReq{
		Search: query.Text,
		Page:   int64(query.Pagination.GetPage()),
		Size:   int64(query.Pagination.GetSize()),
	})
	if err != nil {
		return nil, err
	}
	return dto.GroupsListResponseFromGrpc(res), nil
}
