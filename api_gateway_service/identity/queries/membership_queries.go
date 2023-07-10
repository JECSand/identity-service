package queries

import (
	"context"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	membershipQueryService "github.com/JECSand/identity-service/query_service/protos/membership_query"
	"github.com/opentracing/opentracing-go"
)

type GetMembershipByIdHandler interface {
	Handle(ctx context.Context, query *GetMembershipByIdQuery) (*dto.MembershipResponse, error)
}

type getMembershipByIdHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient membershipQueryService.MembershipQueryServiceClient
}

func NewGetMembershipByIdHandler(log logging.Logger, cfg *config.Config, rsClient membershipQueryService.MembershipQueryServiceClient) *getMembershipByIdHandler {
	return &getMembershipByIdHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (q *getMembershipByIdHandler) Handle(ctx context.Context, query *GetMembershipByIdQuery) (*dto.MembershipResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getMembershipByIdHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := q.rsClient.GetMembershipById(ctx, &membershipQueryService.GetMembershipByIdReq{ID: query.ID.String()})
	if err != nil {
		return nil, err
	}
	return dto.MembershipResponseFromGrpc(res.GetMembership()), nil
}

// GetUserMembershipByGroupIdHandler ...
type GetUserMembershipByGroupIdHandler interface {
	Handle(ctx context.Context, query *GetUserMembershipByGroupIdQuery) (*dto.UserMembershipsListResponse, error)
}

type getUserMembershipByGroupIdHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient membershipQueryService.MembershipQueryServiceClient
}

func NewGetUserMembershipByGroupIHandler(log logging.Logger, cfg *config.Config, rsClient membershipQueryService.MembershipQueryServiceClient) *getUserMembershipByGroupIdHandler {
	return &getUserMembershipByGroupIdHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (s *getUserMembershipByGroupIdHandler) Handle(ctx context.Context, query *GetUserMembershipByGroupIdQuery) (*dto.UserMembershipsListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getUserMembershipByGroupIdHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := s.rsClient.GetUserMembership(ctx, &membershipQueryService.GetUserMembershipReq{
		GroupID: query.GroupID.String(),
		Page:    int64(query.Pagination.GetPage()),
		Size:    int64(query.Pagination.GetSize()),
	})
	if err != nil {
		return nil, err
	}
	return dto.UserMembershipListResponseFromGrpc(res), nil
}

// GetGroupMembershipByUserIdHandler ...
type GetGroupMembershipByUserIdHandler interface {
	Handle(ctx context.Context, query *GetGroupMembershipByUserIdQuery) (*dto.GroupMembershipsListResponse, error)
}

type getGroupMembershipByUserIdHandler struct {
	log      logging.Logger
	cfg      *config.Config
	rsClient membershipQueryService.MembershipQueryServiceClient
}

func NewGetGroupMembershipByUserIdHandler(log logging.Logger, cfg *config.Config, rsClient membershipQueryService.MembershipQueryServiceClient) *getGroupMembershipByUserIdHandler {
	return &getGroupMembershipByUserIdHandler{
		log:      log,
		cfg:      cfg,
		rsClient: rsClient,
	}
}

func (s *getGroupMembershipByUserIdHandler) Handle(ctx context.Context, query *GetGroupMembershipByUserIdQuery) (*dto.GroupMembershipsListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getGroupMembershipByUserIdHandler.Handle")
	defer span.Finish()
	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := s.rsClient.GetGroupMembership(ctx, &membershipQueryService.GetGroupMembershipReq{
		UserID: query.UserID.String(),
		Page:   int64(query.Pagination.GetPage()),
		Size:   int64(query.Pagination.GetSize()),
	})
	if err != nil {
		return nil, err
	}
	return dto.GroupMembershipListResponseFromGrpc(res), nil
}
