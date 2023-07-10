package grpc

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	"github.com/JECSand/identity-service/command_service/identity/metrics"
	"github.com/JECSand/identity-service/command_service/identity/queries"
	"github.com/JECSand/identity-service/command_service/identity/services"
	"github.com/JECSand/identity-service/command_service/mappings"
	membershipCommandService "github.com/JECSand/identity-service/command_service/protos/membership_command"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type membershipGrpcService struct {
	log               logging.Logger
	cfg               *config.Config
	v                 *validator.Validate
	membershipService *services.MembershipService
	metrics           *metrics.CommandServiceMetrics
}

func NewMembershipCommandGrpcService(log logging.Logger, cfg *config.Config, v *validator.Validate, ms *services.MembershipService, metrics *metrics.CommandServiceMetrics) *membershipGrpcService {
	return &membershipGrpcService{
		log:               log,
		cfg:               cfg,
		v:                 v,
		membershipService: ms,
		metrics:           metrics,
	}
}

func (s *membershipGrpcService) CreateMembership(ctx context.Context, req *membershipCommandService.CreateMembershipReq) (*membershipCommandService.CreateMembershipRes, error) {
	s.metrics.CreateMembershipGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "membershipGrpcService.CreateMembership")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	userId, err := uuid.FromString(req.GetUserID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	groupId, err := uuid.FromString(req.GetGroupID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	command := commands.NewCreateMembershipCommand(id, userId, groupId, enums.MembershipStatus(req.GetStatus()), enums.Role(req.GetRole()))
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	err = s.membershipService.Commands.CreateMembership.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("CreateMembership.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &membershipCommandService.CreateMembershipRes{ID: id.String()}, nil
}

func (s *membershipGrpcService) UpdateMembership(ctx context.Context, req *membershipCommandService.UpdateMembershipReq) (*membershipCommandService.UpdateMembershipRes, error) {
	s.metrics.UpdateMembershipGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "membershipGrpcService.UpdateMembership")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	command := commands.NewUpdateMembershipCommand(id, enums.MembershipStatus(req.GetStatus()), enums.Role(req.GetRole()))
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	err = s.membershipService.Commands.UpdateMembership.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("UpdateMembership.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &membershipCommandService.UpdateMembershipRes{}, nil
}

func (s *membershipGrpcService) GetMembershipById(ctx context.Context, req *membershipCommandService.GetMembershipByIdReq) (*membershipCommandService.GetMembershipByIdRes, error) {
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
	found, err := s.membershipService.Queries.GetMembershipById.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetMembershipById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &membershipCommandService.GetMembershipByIdRes{Membership: mappings.CommandMembershipToGrpc(found)}, nil
}

func (s *membershipGrpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
