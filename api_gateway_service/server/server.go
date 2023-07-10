package server

import (
	"context"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/client"
	"github.com/JECSand/identity-service/api_gateway_service/identity/controllers/http/v1"
	"github.com/JECSand/identity-service/api_gateway_service/identity/metrics"
	"github.com/JECSand/identity-service/api_gateway_service/identity/middlewares"
	"github.com/JECSand/identity-service/api_gateway_service/identity/services"
	"github.com/JECSand/identity-service/pkg/authentication"
	"github.com/JECSand/identity-service/pkg/interceptors"
	"github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	authQueryService "github.com/JECSand/identity-service/query_service/protos/auth_query"
	groupQueryService "github.com/JECSand/identity-service/query_service/protos/group_query"
	membershipQueryService "github.com/JECSand/identity-service/query_service/protos/membership_query"
	queryService "github.com/JECSand/identity-service/query_service/protos/user_query"
	"github.com/go-playground/validator"
	"github.com/heptiolabs/healthcheck"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	stackSize      = 1 << 10 // 1 KB
	bodyLimit      = "2M"
	readTimeout    = 15 * time.Second
	writeTimeout   = 15 * time.Second
	gzipLevel      = 5
)

type server struct {
	log  logging.Logger
	auth authentication.Authenticator
	cfg  *config.Config
	v    *validator.Validate
	mw   middlewares.MiddlewareManager
	im   interceptors.InterceptorManager
	echo *echo.Echo
	ps   *services.UserService
	gs   *services.GroupService
	ms   *services.MembershipService
	as   *services.AuthService
	m    *metrics.ApiGatewayMetrics
}

func NewServer(log logging.Logger, auth authentication.Authenticator, cfg *config.Config) *server {
	return &server{
		log:  log,
		auth: auth,
		cfg:  cfg,
		echo: echo.New(), v: validator.New(),
	}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	s.im = interceptors.NewInterceptorManager(s.log, s.auth)
	s.m = metrics.NewApiGatewayMetrics(s.cfg)
	queryServiceClient, err := client.NewQueryServiceClient(ctx, s.cfg, s.im)
	if err != nil {
		return err
	}
	defer queryServiceClient.Close() // nolint: errCheck
	rsClient := queryService.NewQueryServiceClient(queryServiceClient)
	authQueryServiceClient, err := client.NewQueryServiceClient(ctx, s.cfg, s.im)
	if err != nil {
		return err
	}
	defer authQueryServiceClient.Close() // nolint: errCheck
	rsAuthClient := authQueryService.NewAuthQueryServiceClient(authQueryServiceClient)
	groupQueryServiceClient, err := client.NewQueryServiceClient(ctx, s.cfg, s.im)
	if err != nil {
		return err
	}
	defer groupQueryServiceClient.Close() // nolint: errCheck
	rsGroupClient := groupQueryService.NewGroupQueryServiceClient(groupQueryServiceClient)
	membershipQueryServiceClient, err := client.NewQueryServiceClient(ctx, s.cfg, s.im)
	if err != nil {
		return err
	}
	defer membershipQueryServiceClient.Close() // nolint: errCheck
	rsMembershipClient := membershipQueryService.NewMembershipQueryServiceClient(membershipQueryServiceClient)
	kafkaProducer := kafka.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close() // nolint: errCheck
	s.ps = services.NewUserService(s.log, s.cfg, kafkaProducer, rsClient)
	s.gs = services.NewGroupService(s.log, s.cfg, kafkaProducer, rsGroupClient)
	s.ms = services.NewMembershipService(s.log, s.cfg, kafkaProducer, rsMembershipClient)
	s.as = services.NewAuthService(s.log, s.cfg, kafkaProducer, rsAuthClient)
	s.mw = middlewares.NewMiddlewareManager(s.log, s.auth, s.cfg, s.as)
	userHandlers := v1.NewUsersHandlers(s.echo.Group(s.cfg.Http.UsersPath), s.log, s.mw, s.cfg, s.ps, s.ms, s.v, s.m)
	userHandlers.MapRoutes()
	groupHandlers := v1.NewGroupsHandlers(s.echo.Group(s.cfg.Http.GroupsPath), s.log, s.mw, s.cfg, s.gs, s.ms, s.v, s.m)
	groupHandlers.MapRoutes()
	membershipHandlers := v1.NewMembershipsHandlers(s.echo.Group(s.cfg.Http.GroupsPath), s.log, s.mw, s.cfg, s.ms, s.v, s.m)
	membershipHandlers.MapRoutes()
	authHandlers := v1.NewAuthHandlers(s.echo.Group(s.cfg.Http.AuthPath), s.log, s.auth, s.mw, s.cfg, s.as, s.v, s.m)
	authHandlers.MapRoutes()
	go func() {
		if err = s.runHttpServer(); err != nil {
			s.log.Errorf(" s.runHttpServer: %v", err)
			cancel()
		}
	}()
	s.log.Infof("API Gateway Service is listening on PORT: %s", s.cfg.Http.Port)
	s.runMetrics(cancel)
	s.runHealthCheck(ctx)
	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errCheck
		opentracing.SetGlobalTracer(tracer)
	}
	<-ctx.Done()
	if err = s.echo.Server.Shutdown(ctx); err != nil {
		s.log.WarnMsg("echo.Server.Shutdown", err)
	}
	return nil
}

func (s *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()
	health.AddReadinessCheck(s.cfg.ServiceName, healthcheck.AsyncWithContext(ctx, func() error {
		if s.cfg != nil {
			return nil
		}
		return errors.New("Config not loaded")
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
	go func() {
		s.log.Infof("API_Gateway Kubernetes probes listening on port: %s", s.cfg.Probes.Port)
		if err := http.ListenAndServe(s.cfg.Probes.Port, health); err != nil {
			s.log.WarnMsg("ListenAndServe", err)
		}
	}()
}

func (s *server) runMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	go func() {
		metricsServer.GET(s.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: %s", s.cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(s.cfg.Probes.PrometheusPort); err != nil {
			s.log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}

func (s *server) runHttpServer() error {
	s.mapRoutes()
	s.echo.Server.ReadTimeout = readTimeout
	s.echo.Server.WriteTimeout = writeTimeout
	s.echo.Server.MaxHeaderBytes = maxHeaderBytes
	return s.echo.Start(s.cfg.Http.Port)
}

func (s *server) mapRoutes() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "API Gateway"
	docs.SwaggerInfo.Description = "API Gateway CQRS microservices."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.Use(s.mw.RequestLoggerMiddleware)
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	s.echo.Use(middleware.BodyLimit(bodyLimit))
}
