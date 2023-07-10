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

type CreateMembershipCmdHandler interface {
	Handle(ctx context.Context, command *CreateMembershipCommand) error
}

type createMembershipHandler struct {
	log           logging.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewCreateMembershipHandler(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *createMembershipHandler {
	return &createMembershipHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (c *createMembershipHandler) Handle(ctx context.Context, command *CreateMembershipCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createMembershipHandler.Handle")
	defer span.Finish()
	createDTO := &kafkaMessages.MembershipCreate{
		ID:      command.CreateDto.ID.String(),
		UserID:  command.CreateDto.UserID.String(),
		GroupID: command.CreateDto.GroupID.String(),
		Status:  int64(command.CreateDto.Status),
		Role:    int64(command.CreateDto.Role),
	}
	dtoBytes, err := proto.Marshal(createDTO)
	if err != nil {
		return err
	}
	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.MembershipCreate.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}

// UpdateMembershipCmdHandler ...
type UpdateMembershipCmdHandler interface {
	Handle(ctx context.Context, command *UpdateMembershipCommand) error
}

type updateMembershipCmdHandler struct {
	log           logging.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewUpdateMembershipHandler(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *updateMembershipCmdHandler {
	return &updateMembershipCmdHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (c *updateMembershipCmdHandler) Handle(ctx context.Context, command *UpdateMembershipCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateMembershipCmdHandler.Handle")
	defer span.Finish()
	updateDTO := &kafkaMessages.MembershipUpdate{
		ID:     command.UpdateDto.ID.String(),
		Status: int64(command.UpdateDto.Status),
		Role:   int64(command.UpdateDto.Role),
	}
	dtoBytes, err := proto.Marshal(updateDTO)
	if err != nil {
		return err
	}
	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.MembershipUpdate.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}

// DeleteMembershipCmdHandler ...
type DeleteMembershipCmdHandler interface {
	Handle(ctx context.Context, command *DeleteMembershipCommand) error
}

type deleteMembershipHandler struct {
	log           logging.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewDeleteMembershipHandler(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *deleteMembershipHandler {
	return &deleteMembershipHandler{log: log, cfg: cfg, kafkaProducer: kafkaProducer}
}

func (c *deleteMembershipHandler) Handle(ctx context.Context, command *DeleteMembershipCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteMembershipHandler.Handle")
	defer span.Finish()
	deleteDTO := &kafkaMessages.MembershipDelete{ID: command.ID.String()}
	dtoBytes, err := proto.Marshal(deleteDTO)
	if err != nil {
		return err
	}
	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.MembershipDelete.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}
