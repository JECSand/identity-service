package grpc

import (
	"context"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/JECSand/identity-service/query_service/identity/events"
	"github.com/JECSand/identity-service/query_service/identity/metrics"
	"github.com/JECSand/identity-service/query_service/identity/queries"
	"github.com/JECSand/identity-service/query_service/identity/services"
	authQueryService "github.com/JECSand/identity-service/query_service/protos/auth_query"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type authGrpcService struct {
	log     logging.Logger
	cfg     *config.Config
	v       *validator.Validate
	us      *services.UserService
	as      *services.AuthService
	metrics *metrics.QueryServiceMetrics
}

func NewAuthQueryGrpcService(
	log logging.Logger,
	cfg *config.Config,
	v *validator.Validate,
	us *services.UserService,
	as *services.AuthService,
	metrics *metrics.QueryServiceMetrics,
) *authGrpcService {
	return &authGrpcService{
		log:     log,
		cfg:     cfg,
		v:       v,
		us:      us,
		as:      as,
		metrics: metrics,
	}
}

func (s *authGrpcService) UpdatePassword(ctx context.Context, req *authQueryService.PasswordUpdateReq) (*authQueryService.PasswordUpdateRes, error) {
	s.metrics.UpdatePasswordGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "authGrpcService.UpdatePassword")
	defer span.Finish()
	event := events.NewUpdatePasswordEvent(req.GetID(), req.GetNewPassword(), time.Now())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err := s.as.Events.UpdatePassword.Handle(ctx, event); err != nil {
		s.log.WarnMsg("UpdatePassword.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &authQueryService.PasswordUpdateRes{Status: 200}, nil
}

func (s *authGrpcService) BlacklistToken(ctx context.Context, req *authQueryService.BlacklistTokenReq) (*authQueryService.BlacklistTokenRes, error) {
	s.metrics.BlacklistTokenGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "authGrpcService.BlacklistToken")
	defer span.Finish()
	event := events.NewBlacklistTokenEvent(req.GetID(), req.GetAccessToken(), time.Now(), time.Now())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err := s.as.Events.BlacklistToken.Handle(ctx, event); err != nil {
		s.log.WarnMsg("BlacklistToken.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &authQueryService.BlacklistTokenRes{ID: event.ID, Status: 200}, nil
}

func (s *authGrpcService) Authenticate(ctx context.Context, req *authQueryService.AuthenticateReq) (*authQueryService.AuthenticateRes, error) {
	s.metrics.AuthenticateGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "authGrpcService.Authenticate")
	defer span.Finish()
	query := queries.NewAuthenticateQuery(req.GetEmail(), req.GetPassword())
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	user, err := s.as.Queries.Authenticate.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("Authenticate.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &authQueryService.AuthenticateRes{User: entities.AuthUserToGrpcMessage(user), Status: 200}, nil
}

func (s *authGrpcService) Validate(ctx context.Context, req *authQueryService.ValidateReq) (*authQueryService.ValidateRes, error) {
	s.metrics.ValidateGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "authGrpcService.Validate")
	defer span.Finish()
	id, err := uuid.FromString(req.GetUserID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	query := queries.NewValidateQuery(id, req.GetAccessToken(), enums.ValidationType(req.GetValidationType()))
	if err = s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	user, err := s.as.Queries.Validate.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("Validate.Handle", err)
		return nil, s.errResponse(codes.Unauthenticated, err)
	}
	s.metrics.SuccessGrpcRequests.Inc()
	return &authQueryService.ValidateRes{User: entities.AuthUserToGrpcMessage(user), Status: 200}, nil
}

func (s *authGrpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
