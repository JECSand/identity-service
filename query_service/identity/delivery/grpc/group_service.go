package grpc

import (
	"context"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/JECSand/identity-service/query_service/identity/events"
	"github.com/JECSand/identity-service/query_service/identity/metrics"
	"github.com/JECSand/identity-service/query_service/identity/queries"
	"github.com/JECSand/identity-service/query_service/identity/services"
	groupQueryService "github.com/JECSand/identity-service/query_service/protos/group_query"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type groupGrpcService struct {
	log     logging.Logger
	cfg     *config.Config
	v       *validator.Validate
	gs      *services.GroupService
	metrics *metrics.QueryServiceMetrics
}

func NewGroupQueryGrpcService(
	log logging.Logger,
	cfg *config.Config,
	v *validator.Validate,
	gs *services.GroupService,
	metrics *metrics.QueryServiceMetrics,
) *groupGrpcService {
	return &groupGrpcService{
		log:     log,
		cfg:     cfg,
		v:       v,
		gs:      gs,
		metrics: metrics,
	}
}

func (s *groupGrpcService) CreateGroup(ctx context.Context, req *groupQueryService.CreateGroupReq) (*groupQueryService.CreateGroupRes, error) {
	s.metrics.CreateGroupGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "groupGrpcService.CreateGroup")
	defer span.Finish()
	// TODO - ADD LOGIC FOR ROOT AND ACTIVE BELOW
	event := events.NewCreateGroupEvent(req.GetID(), req.GetName(), req.GetDescription(), req.GetCreatorID(), false, time.Now(), time.Now())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err := s.gs.Events.CreateGroup.Handle(ctx, event); err != nil {
		s.log.WarnMsg("CreateGroup.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &groupQueryService.CreateGroupRes{ID: req.GetID()}, nil
}

func (s *groupGrpcService) UpdateGroup(ctx context.Context, req *groupQueryService.UpdateGroupReq) (*groupQueryService.UpdateGroupRes, error) {
	s.metrics.UpdateGroupGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "groupGrpcService.UpdateGroup")
	defer span.Finish()
	command := events.NewUpdateGroupEvent(req.GetID(), req.GetName(), req.GetDescription(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err := s.gs.Events.UpdateGroup.Handle(ctx, command); err != nil {
		s.log.WarnMsg("UpdateGroup.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &groupQueryService.UpdateGroupRes{ID: req.GetID()}, nil
}

func (s *groupGrpcService) GetGroupById(ctx context.Context, req *groupQueryService.GetGroupByIdReq) (*groupQueryService.GetGroupByIdRes, error) {
	s.metrics.GetGroupByIdGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "groupGrpcService.GetGroupById")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	query := queries.NewGetGroupByIdQuery(id)
	if err = s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	group, err := s.gs.Queries.GetGroupById.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetGroupById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &groupQueryService.GetGroupByIdRes{Group: entities.GroupToGrpcMessage(group)}, nil
}

func (s *groupGrpcService) SearchGroup(ctx context.Context, req *groupQueryService.SearchGroupReq) (*groupQueryService.SearchGroupRes, error) {
	s.metrics.SearchGroupGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "groupGrpcService.SearchGroup")
	defer span.Finish()
	pq := utilities.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))
	query := queries.NewSearchGroupQuery(req.GetSearch(), pq)
	groupsList, err := s.gs.Queries.SearchGroup.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("SearchGroup.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return entities.GroupListToGrpc(groupsList), nil
}

func (s *groupGrpcService) DeleteGroupByID(ctx context.Context, req *groupQueryService.DeleteGroupByIdReq) (*groupQueryService.DeleteGroupByIdRes, error) {
	s.metrics.DeleteMembershipGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "groupGrpcService.DeleteGroupByID")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err = s.gs.Events.DeleteGroup.Handle(ctx, events.NewDeleteGroupEvent(id)); err != nil {
		s.log.WarnMsg("DeleteGroup.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &groupQueryService.DeleteGroupByIdRes{}, nil
}

func (s *groupGrpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
