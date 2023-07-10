package commands

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/models"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	"github.com/JECSand/identity-service/command_service/mappings"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	kafkaMessages "github.com/JECSand/identity-service/protos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// BlacklistTokenCmdHandler ...
type BlacklistTokenCmdHandler interface {
	Handle(ctx context.Context, command *BlacklistTokenCommand) error
}

type blacklistTokenHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewBlacklistTokenHandler ...
func NewBlacklistTokenHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *blacklistTokenHandler {
	return &blacklistTokenHandler{
		log:           log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *blacklistTokenHandler) Handle(ctx context.Context, command *BlacklistTokenCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blacklistTokenHandler.Handle")
	defer span.Finish()
	blDTO := &models.Blacklist{
		ID:          command.ID,
		AccessToken: command.AccessToken,
	}
	bl, err := c.pgRepo.BlacklistToken(ctx, blDTO)
	if err != nil {
		return err
	}
	msg := &kafkaMessages.TokenBlacklisted{Blacklist: mappings.BlacklistToGrpcMessage(bl)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.TokenBlacklisted.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}

// PasswordUpdateCmdHandler ...
type PasswordUpdateCmdHandler interface {
	Handle(ctx context.Context, command *PasswordUpdateCommand) error
}

type passwordUpdateHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewUpdatePasswordHandler ...
func NewUpdatePasswordHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *passwordUpdateHandler {
	return &passwordUpdateHandler{
		log:           log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *passwordUpdateHandler) Handle(ctx context.Context, command *PasswordUpdateCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "passwordUpdateHandler.Handle")
	defer span.Finish()
	authDTO := &models.User{
		ID:       command.ID,
		Password: command.NewPassword,
	}
	// TODO: QUERY FOR USER RECORD BY ID, HASH COMMAND PASSWORD, AND COMPARE CURRENT COMMAND PASSWORD WITH USER PASSWORD
	if err := authDTO.HashPassword(); err != nil {
		return err
	}
	user, err := c.pgRepo.UpdateUserPassword(ctx, authDTO)
	if err != nil {
		return err
	}
	msg := &kafkaMessages.PasswordUpdated{
		ID:          user.ID.String(),
		NewPassword: authDTO.Password,
		Status:      200,
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.PasswordUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}
