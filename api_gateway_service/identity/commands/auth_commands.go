package commands

import (
	"context"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	kafkaMessages "github.com/JECSand/identity-service/protos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type BlacklistTokenCmdHandler interface {
	Handle(ctx context.Context, command *BlacklistTokenCommand) error
}

type blacklistTokenHandler struct {
	log           logging.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewBlacklistTokenHandler(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *blacklistTokenHandler {
	return &blacklistTokenHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (c *blacklistTokenHandler) Handle(ctx context.Context, command *BlacklistTokenCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blacklistTokenHandler.Handle")
	defer span.Finish()
	blacklistDTO := &kafkaMessages.TokenBlacklist{
		ID:          command.BlacklistDto.ID.String(),
		AccessToken: command.BlacklistDto.AccessToken,
	}
	dtoBytes, err := proto.Marshal(blacklistDTO)
	if err != nil {
		return err
	}
	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.TokenBlacklist.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}

// UpdatePasswordCmdHandler ...
type UpdatePasswordCmdHandler interface {
	Handle(ctx context.Context, command *UpdatePasswordCommand) error
}

type updatePasswordCmdHandler struct {
	log           logging.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewUpdatePasswordHandler(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *updatePasswordCmdHandler {
	return &updatePasswordCmdHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (c *updatePasswordCmdHandler) Handle(ctx context.Context, command *UpdatePasswordCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updatePasswordCmdHandler.Handle")
	defer span.Finish()
	updateDTO := &kafkaMessages.PasswordUpdate{
		ID:              command.UpdateDto.ID.String(),
		CurrentPassword: command.UpdateDto.CurrentPassword,
		NewPassword:     command.UpdateDto.NewPassword,
	}
	dtoBytes, err := proto.Marshal(updateDTO)
	if err != nil {
		return err
	}
	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.PasswordUpdate.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}
