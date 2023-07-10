package server

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	grpc3 "github.com/JECSand/identity-service/command_service/identity/delivery/grpc"
	kafkaConsumer "github.com/JECSand/identity-service/command_service/identity/delivery/kafka"
	"github.com/JECSand/identity-service/command_service/identity/metrics"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	"github.com/JECSand/identity-service/command_service/identity/services"
	authCommandService "github.com/JECSand/identity-service/command_service/protos/auth_command"
	groupCommandService "github.com/JECSand/identity-service/command_service/protos/group_command"
	membershipCommandService "github.com/JECSand/identity-service/command_service/protos/membership_command"
	commandService "github.com/JECSand/identity-service/command_service/protos/user_command"
	"github.com/JECSand/identity-service/pkg/authentication"
	"github.com/JECSand/identity-service/pkg/constants"
	"github.com/JECSand/identity-service/pkg/interceptors"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/postgres"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/go-playground/validator"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
	stackSize         = 1 << 10 // 1 KB
)

type server struct {
	log               logging.Logger
	auth              authentication.Authenticator
	cfg               *config.Config
	v                 *validator.Validate
	kafkaConn         *kafka.Conn
	userService       *services.UserService
	groupService      *services.GroupService
	membershipService *services.MembershipService
	authService       *services.AuthService
	im                interceptors.InterceptorManager
	pgConn            *pgxpool.Pool
	metrics           *metrics.CommandServiceMetrics
}

func NewServer(log logging.Logger, cfg *config.Config) *server {
	return &server{
		log: log,
		cfg: cfg,
		v:   validator.New(),
	}
}

func (s *server) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, s.cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}
	s.kafkaConn = kafkaConn
	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}
	s.log.Infof("kafka connected to brokers: %+v", brokers)
	return nil
}

func (s *server) initKafkaTopics(ctx context.Context) {
	controller, err := s.kafkaConn.Controller()
	if err != nil {
		s.log.WarnMsg("kafkaConn.Controller", err)
		return
	}
	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.log.Infof("kafka controller uri: %s", controllerURI)
	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errCheck
	s.log.Infof("established new kafka controller connection: %s", controllerURI)
	userCreateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.UserCreate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.UserCreate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.UserCreate.ReplicationFactor,
	}
	userCreatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.UserCreated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.UserCreated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.UserCreated.ReplicationFactor,
	}
	userUpdateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.UserUpdate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.UserUpdate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.UserUpdate.ReplicationFactor,
	}
	userUpdatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.UserUpdated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.UserUpdated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.UserUpdated.ReplicationFactor,
	}
	userDeleteTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.UserDelete.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.UserDelete.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.UserDelete.ReplicationFactor,
	}
	userDeletedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.UserDeleted.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.UserDeleted.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.UserDeleted.ReplicationFactor,
	}
	groupCreateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.GroupCreate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.GroupCreate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.GroupCreate.ReplicationFactor,
	}
	groupCreatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.GroupCreated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.GroupCreated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.GroupCreated.ReplicationFactor,
	}
	groupUpdateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.GroupUpdate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.GroupUpdate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.GroupUpdate.ReplicationFactor,
	}
	groupUpdatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.GroupUpdated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.GroupUpdated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.GroupUpdated.ReplicationFactor,
	}
	groupDeleteTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.GroupDelete.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.GroupDelete.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.GroupDelete.ReplicationFactor,
	}
	groupDeletedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.GroupDeleted.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.GroupDeleted.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.GroupDeleted.ReplicationFactor,
	}
	membershipCreateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.MembershipCreate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.MembershipCreate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.MembershipCreate.ReplicationFactor,
	}
	membershipCreatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.MembershipCreated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.MembershipCreated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.MembershipCreated.ReplicationFactor,
	}
	membershipUpdateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.MembershipUpdate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.MembershipUpdate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.MembershipUpdate.ReplicationFactor,
	}
	membershipUpdatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.MembershipUpdated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.MembershipUpdated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.MembershipUpdated.ReplicationFactor,
	}
	membershipDeleteTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.MembershipDelete.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.MembershipDelete.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.MembershipDelete.ReplicationFactor,
	}
	membershipDeletedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.MembershipDeleted.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.MembershipDeleted.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.MembershipDeleted.ReplicationFactor,
	}
	tokenBlacklistTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.TokenBlacklist.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.TokenBlacklist.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.TokenBlacklist.ReplicationFactor,
	}
	tokenBlacklistedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.TokenBlacklisted.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.TokenBlacklisted.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.TokenBlacklisted.ReplicationFactor,
	}
	passwordUpdateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.PasswordUpdate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.PasswordUpdate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.PasswordUpdate.ReplicationFactor,
	}
	passwordUpdatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.PasswordUpdated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.PasswordUpdated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.PasswordUpdated.ReplicationFactor,
	}
	if err = conn.CreateTopics(
		userCreateTopic,
		userUpdateTopic,
		userCreatedTopic,
		userUpdatedTopic,
		userDeleteTopic,
		userDeletedTopic,
		groupCreateTopic,
		groupUpdateTopic,
		groupCreatedTopic,
		groupUpdatedTopic,
		groupDeleteTopic,
		groupDeletedTopic,
		membershipCreateTopic,
		membershipUpdateTopic,
		membershipCreatedTopic,
		membershipUpdatedTopic,
		membershipDeleteTopic,
		membershipDeletedTopic,
		tokenBlacklistTopic,
		tokenBlacklistedTopic,
		passwordUpdateTopic,
		passwordUpdatedTopic,
	); err != nil {
		s.log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}
	s.log.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{
		userCreateTopic,
		userUpdateTopic,
		userCreatedTopic,
		userUpdatedTopic,
		userDeleteTopic,
		userDeletedTopic,
		groupCreateTopic,
		groupUpdateTopic,
		groupCreatedTopic,
		groupUpdatedTopic,
		groupDeleteTopic,
		groupDeletedTopic,
		membershipCreateTopic,
		membershipUpdateTopic,
		membershipCreatedTopic,
		membershipUpdatedTopic,
		membershipDeleteTopic,
		membershipDeletedTopic,
		tokenBlacklistTopic,
		tokenBlacklistedTopic,
		passwordUpdateTopic,
		passwordUpdatedTopic,
	})
}

func (s *server) getConsumerGroupTopics() []string {
	return []string{
		s.cfg.KafkaTopics.UserCreate.TopicName,
		s.cfg.KafkaTopics.UserUpdate.TopicName,
		s.cfg.KafkaTopics.UserDelete.TopicName,
		s.cfg.KafkaTopics.GroupCreate.TopicName,
		s.cfg.KafkaTopics.GroupUpdate.TopicName,
		s.cfg.KafkaTopics.GroupDelete.TopicName,
		s.cfg.KafkaTopics.MembershipCreate.TopicName,
		s.cfg.KafkaTopics.MembershipUpdate.TopicName,
		s.cfg.KafkaTopics.MembershipDelete.TopicName,
		s.cfg.KafkaTopics.TokenBlacklist.TopicName,
		s.cfg.KafkaTopics.PasswordUpdate.TopicName,
	}
}

func (s *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()
	health.AddLivenessCheck(s.cfg.ServiceName, healthcheck.AsyncWithContext(ctx, func() error {
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
	health.AddReadinessCheck(constants.Postgres, healthcheck.AsyncWithContext(ctx, func() error {
		return s.pgConn.Ping(ctx)
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
	health.AddReadinessCheck(constants.Kafka, healthcheck.AsyncWithContext(ctx, func() error {
		_, err := s.kafkaConn.Brokers()
		if err != nil {
			return err
		}
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
	go func() {
		s.log.Infof("Command service Kubernetes probes listening on port: %s", s.cfg.Probes.Port)
		if err := http.ListenAndServe(s.cfg.Probes.Port, health); err != nil {
			s.log.WarnMsg("ListenAndServe", err)
		}
	}()
}

func (s *server) runMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	go func() {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         stackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(s.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: %s", s.cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(s.cfg.Probes.PrometheusPort); err != nil {
			s.log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}

func (s *server) runInitializations(ctx context.Context) {
	count, err := s.userService.Queries.CountUsers.Handle(ctx)
	if err != nil {
		s.log.Errorf("runInitializations: %v", err)
	}
	if count == 0 {
		r := s.cfg.Initialization.Users.Root
		id, err := utilities.NewID()
		if err != nil {
			s.log.WarnMsg("utilities.NewID", err)
		}
		command := commands.NewCreateUserCommand(id, r.Email, r.Username, r.Password, true, true)
		if err = s.v.StructCtx(ctx, command); err != nil {
			s.log.WarnMsg("validate", err)
		}
		err = s.userService.Commands.CreateUser.Handle(ctx, command)
		if err != nil {
			s.log.WarnMsg("CreateUser.Handle", err)
		}
	}
}

func (s *server) newCommandGrpcServer() (func() error, *grpc.Server, error) {
	l, err := net.Listen("tcp", s.cfg.GRPC.Port)
	if err != nil {
		return nil, nil, errors.Wrap(err, "net.Listen")
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			s.im.Logger,
		)),
	)
	commandGrpcWriter := grpc3.NewCommandGrpcService(s.log, s.cfg, s.v, s.userService, s.authService, s.metrics)
	commandService.RegisterCommandServiceServer(grpcServer, commandGrpcWriter)
	grpc_prometheus.Register(grpcServer)
	authCommandGrpcWriter := grpc3.NewAuthCommandGrpcService(s.log, s.cfg, s.v, s.userService, s.authService, s.metrics)
	authCommandService.RegisterAuthCommandServiceServer(grpcServer, authCommandGrpcWriter)
	groupCommandGrpcWriter := grpc3.NewGroupCommandGrpcService(s.log, s.cfg, s.v, s.groupService, s.metrics)
	groupCommandService.RegisterGroupCommandServiceServer(grpcServer, groupCommandGrpcWriter)
	membershipCommandGrpcWriter := grpc3.NewMembershipCommandGrpcService(s.log, s.cfg, s.v, s.membershipService, s.metrics)
	membershipCommandService.RegisterMembershipCommandServiceServer(grpcServer, membershipCommandGrpcWriter)
	grpc_prometheus.Register(grpcServer)
	if s.cfg.GRPC.Development {
		reflection.Register(grpcServer)
	}
	go func() {
		s.log.Infof("Command gRPC server is listening on port: %s", s.cfg.GRPC.Port)
		s.log.Fatal(grpcServer.Serve(l))
	}()
	return l.Close, grpcServer, nil
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	s.im = interceptors.NewInterceptorManager(s.log, s.auth)
	s.metrics = metrics.NewCommandServiceMetrics(s.cfg)
	pgxConn, err := postgres.NewPostgresConn(s.cfg.Postgresql)
	if err != nil {
		return errors.Wrap(err, "postgresql.NewPostgresConn")
	}
	s.pgConn = pgxConn
	s.log.Infof("postgres connected: %v", pgxConn.Stat().TotalConns())
	defer pgxConn.Close()
	kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close() // nolint: errCheck
	repo := repositories.NewRepository(s.log, s.cfg, pgxConn)
	s.userService = services.NewUserService(s.log, s.cfg, repo, kafkaProducer)
	s.groupService = services.NewGroupService(s.log, s.cfg, repo, kafkaProducer)
	s.membershipService = services.NewMembershipService(s.log, s.cfg, repo, kafkaProducer)
	s.authService = services.NewAuthService(s.log, s.cfg, repo, kafkaProducer)
	identityMessageProcessor := kafkaConsumer.NewIdentityMessageProcessor(
		s.log,
		s.cfg,
		s.v,
		s.userService,
		s.groupService,
		s.membershipService,
		s.authService,
		s.metrics,
	)
	s.log.Info("Starting Writer Kafka consumers")
	cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), kafkaConsumer.PoolSize, identityMessageProcessor.ProcessMessages)
	closeGrpcServer, grpcServer, err := s.newCommandGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer() // nolint: errCheck
	if err = s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close() // nolint: errCheck
	if s.cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}
	s.runHealthCheck(ctx)
	s.runMetrics(cancel)
	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errCheck
		opentracing.SetGlobalTracer(tracer)
	}
	s.runInitializations(ctx)
	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
