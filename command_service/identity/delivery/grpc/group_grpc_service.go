package grpc

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	"github.com/JECSand/identity-service/command_service/identity/metrics"
	"github.com/JECSand/identity-service/command_service/identity/queries"
	"github.com/JECSand/identity-service/command_service/identity/services"
	"github.com/JECSand/identity-service/command_service/mappings"
	groupCommandService "github.com/JECSand/identity-service/command_service/protos/group_command"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type groupGrpcService struct {
	log          logging.Logger
	cfg          *config.Config
	v            *validator.Validate
	groupService *services.GroupService
	metrics      *metrics.CommandServiceMetrics
}

func NewGroupCommandGrpcService(log logging.Logger, cfg *config.Config, v *validator.Validate, gs *services.GroupService, metrics *metrics.CommandServiceMetrics) *groupGrpcService {
	return &groupGrpcService{
		log:          log,
		cfg:          cfg,
		v:            v,
		groupService: gs,
		metrics:      metrics,
	}
}

func (s *groupGrpcService) CreateGroup(ctx context.Context, req *groupCommandService.CreateGroupReq) (*groupCommandService.CreateGroupRes, error) {
	s.metrics.CreateGroupGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "groupGrpcService.CreateGroup")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	creatorId, err := uuid.FromString(req.GetCreatorID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	command := commands.NewCreateGroupCommand(id, req.GetName(), req.GetDescription(), creatorId, false)
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	err = s.groupService.Commands.CreateGroup.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("CreateGroup.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &groupCommandService.CreateGroupRes{ID: id.String()}, nil
}

func (s *groupGrpcService) UpdateGroup(ctx context.Context, req *groupCommandService.UpdateGroupReq) (*groupCommandService.UpdateGroupRes, error) {
	s.metrics.UpdateGroupGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "groupGrpcService.UpdateGroup")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	command := commands.NewUpdateGroupCommand(id, req.GetName(), req.GetDescription())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	err = s.groupService.Commands.UpdateGroup.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("UpdateGroup.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &groupCommandService.UpdateGroupRes{}, nil
}

func (s *groupGrpcService) GetGroupById(ctx context.Context, req *groupCommandService.GetGroupByIdReq) (*groupCommandService.GetGroupByIdRes, error) {
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
	found, err := s.groupService.Queries.GetGroupById.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetGroupById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &groupCommandService.GetGroupByIdRes{Group: mappings.CommandGroupToGrpc(found)}, nil
}

func (s *groupGrpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
