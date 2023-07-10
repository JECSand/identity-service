package grpc

import (
	"context"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/JECSand/identity-service/query_service/identity/events"
	"github.com/JECSand/identity-service/query_service/identity/metrics"
	"github.com/JECSand/identity-service/query_service/identity/queries"
	"github.com/JECSand/identity-service/query_service/identity/services"
	membershipQueryService "github.com/JECSand/identity-service/query_service/protos/membership_query"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type membershipGrpcService struct {
	log     logging.Logger
	cfg     *config.Config
	v       *validator.Validate
	ms      *services.MembershipService
	metrics *metrics.QueryServiceMetrics
}

func NewMembershipQueryGrpcService(
	log logging.Logger,
	cfg *config.Config,
	v *validator.Validate,
	ms *services.MembershipService,
	metrics *metrics.QueryServiceMetrics,
) *membershipGrpcService {
	return &membershipGrpcService{
		log:     log,
		cfg:     cfg,
		v:       v,
		ms:      ms,
		metrics: metrics,
	}
}

func (s *membershipGrpcService) CreateMembership(ctx context.Context, req *membershipQueryService.CreateMembershipReq) (*membershipQueryService.CreateMembershipRes, error) {
	s.metrics.CreateMembershipGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "membershipGrpcService.CreateMembership")
	defer span.Finish()
	datetime := time.Now()
	createdMembership := events.NewCreatedMembership(
		req.GetMembership().GetID(),
		req.GetMembership().GetUserID(),
		req.GetMembership().GetGroupID(),
		enums.MembershipStatus(req.GetMembership().GetStatus()),
		enums.Role(req.GetMembership().GetRole()),
		datetime,
		datetime,
	)
	createdUserMembership := events.NewCreatedUserMembership(
		req.GetUserMembership().GetID(),
		req.GetUserMembership().GetGroupID(),
		req.GetUserMembership().GetUserID(),
		req.GetUserMembership().GetMembershipID(),
		req.GetUserMembership().GetEmail(),
		req.GetUserMembership().GetUsername(),
		enums.MembershipStatus(req.GetUserMembership().GetStatus()),
		enums.Role(req.GetUserMembership().GetRole()),
		datetime,
		datetime,
	)
	createdGroupMembership := events.NewCreatedGroupMembership(
		req.GetGroupMembership().GetID(),
		req.GetGroupMembership().GetGroupID(),
		req.GetGroupMembership().GetUserID(),
		req.GetGroupMembership().GetMembershipID(),
		req.GetGroupMembership().GetName(),
		req.GetGroupMembership().GetDescription(),
		enums.MembershipStatus(req.GetGroupMembership().GetStatus()),
		enums.Role(req.GetGroupMembership().GetRole()),
		req.GetGroupMembership().GetCreator(),
		datetime,
		datetime,
	)
	event := events.NewCreateMembershipEvent(createdMembership, createdUserMembership, createdGroupMembership)
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err := s.ms.Events.CreateMembership.Handle(ctx, event); err != nil {
		s.log.WarnMsg("CreateMembership.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &membershipQueryService.CreateMembershipRes{ID: req.GetMembership().GetID()}, nil
}

func (s *membershipGrpcService) UpdateMembership(ctx context.Context, req *membershipQueryService.UpdateMembershipReq) (*membershipQueryService.UpdateMembershipRes, error) {
	s.metrics.UpdateMembershipGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "membershipGrpcService.UpdateMembership")
	defer span.Finish()
	command := events.NewUpdateMembershipEvent(req.GetID(), enums.MembershipStatus(req.GetStatus()), enums.Role(req.GetRole()), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err := s.ms.Events.UpdateMembership.Handle(ctx, command); err != nil {
		s.log.WarnMsg("UpdateMembership.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &membershipQueryService.UpdateMembershipRes{ID: req.GetID()}, nil
}

func (s *membershipGrpcService) GetMembershipById(ctx context.Context, req *membershipQueryService.GetMembershipByIdReq) (*membershipQueryService.GetMembershipByIdRes, error) {
	s.metrics.GetMembershipByIdGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "membershipGrpcService.GetMembershipById")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	query := queries.NewGetMembershipByIdQuery(id)
	if err = s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	membership, err := s.ms.Queries.GetMembershipById.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetMembershipById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &membershipQueryService.GetMembershipByIdRes{Membership: entities.MembershipToGrpcMessage(membership)}, nil
}

func (s *membershipGrpcService) GetGroupMembership(ctx context.Context, req *membershipQueryService.GetGroupMembershipReq) (*membershipQueryService.GetGroupMembershipRes, error) {
	s.metrics.GetGroupMembershipGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "membershipGrpcService.GetGroupMembership")
	defer span.Finish()
	id, err := uuid.FromString(req.GetUserID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	pq := utilities.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))
	query := queries.NewGetGroupMembershipQuery(id, pq)
	groupsList, err := s.ms.Queries.GetGroupMembership.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetGroupMembership.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return entities.GroupMembershipListToGrpc(groupsList), nil
}

func (s *membershipGrpcService) GetUserMembership(ctx context.Context, req *membershipQueryService.GetUserMembershipReq) (*membershipQueryService.GetUserMembershipRes, error) {
	s.metrics.GetUserMembershipGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "membershipGrpcService.GetUserMembership")
	defer span.Finish()
	id, err := uuid.FromString(req.GetGroupID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	pq := utilities.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))
	query := queries.NewGetUserMembershipQuery(id, pq)
	usersList, err := s.ms.Queries.GetUserMembership.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetUserMembership.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return entities.UserMembershipListToGrpc(usersList), nil
}

func (s *membershipGrpcService) DeleteMembershipByID(ctx context.Context, req *membershipQueryService.DeleteMembershipByIdReq) (*membershipQueryService.DeleteMembershipByIdRes, error) {
	s.metrics.DeleteMembershipGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "membershipGrpcService.DeleteMembershipByID")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err = s.ms.Events.DeleteMembership.Handle(ctx, events.NewDeleteMembershipEvent(id)); err != nil {
		s.log.WarnMsg("DeleteMembership.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &membershipQueryService.DeleteMembershipByIdRes{}, nil
}

func (s *membershipGrpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
