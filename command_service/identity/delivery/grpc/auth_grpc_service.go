package grpc

import (
	"context"
	"errors"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	"github.com/JECSand/identity-service/command_service/identity/metrics"
	"github.com/JECSand/identity-service/command_service/identity/queries"
	"github.com/JECSand/identity-service/command_service/identity/services"
	"github.com/JECSand/identity-service/command_service/protos/auth_command"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authGrpcService struct {
	log         logging.Logger
	cfg         *config.Config
	v           *validator.Validate
	userService *services.UserService
	authService *services.AuthService
	metrics     *metrics.CommandServiceMetrics
}

func NewAuthCommandGrpcService(log logging.Logger, cfg *config.Config, v *validator.Validate, us *services.UserService, as *services.AuthService, metrics *metrics.CommandServiceMetrics) *authGrpcService {
	return &authGrpcService{
		log:         log,
		cfg:         cfg,
		v:           v,
		userService: us,
		authService: as,
		metrics:     metrics,
	}
}

func (s *authGrpcService) BlacklistToken(ctx context.Context, req *authCommandService.BlacklistTokenReq) (*authCommandService.BlacklistTokenRes, error) {
	s.metrics.BlacklistTokenGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "authGrpcService.BlacklistToken")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	command := commands.NewBlacklistTokenCommand(id, req.GetAccessToken())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	err = s.authService.Commands.BlacklistToken.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("BlacklistToken.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &authCommandService.BlacklistTokenRes{ID: id.String()}, nil
}

func (s *authGrpcService) UpdatePassword(ctx context.Context, req *authCommandService.UpdatePasswordReq) (*authCommandService.UpdatePasswordRes, error) {
	s.metrics.PasswordUpdateGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "authGrpcService.UpdatePassword")
	defer span.Finish()
	id, err := uuid.FromString(req.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	command := commands.NewUpdatePasswordCommand(id, req.GetCurrentPassword(), req.GetNewPassword())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	err = s.authService.Commands.UpdatePassword.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("UpdatePassword.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &authCommandService.UpdatePasswordRes{Status: 200}, nil
}

func (s *authGrpcService) CheckTokenBlacklist(ctx context.Context, req *authCommandService.CheckBlacklistReq) (*authCommandService.CheckBlacklistRes, error) {
	s.metrics.CheckTokenBlacklistGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "authGrpcService.CheckTokenBlacklist")
	defer span.Finish()
	query := queries.NewCheckTokenBlacklistQuery(req.GetAccessToken())
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	_, err := s.authService.Queries.CheckTokenBlacklist.Handle(ctx, query)
	if err == nil {
		err = errors.New("token is blacklisted")
		s.log.WarnMsg("CheckTokenBlacklist.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &authCommandService.CheckBlacklistRes{Status: 200}, nil
}

func (s *authGrpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
