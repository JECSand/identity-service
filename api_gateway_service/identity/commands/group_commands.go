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

type CreateGroupCmdHandler interface {
	Handle(ctx context.Context, command *CreateGroupCommand) error
}

type createGroupHandler struct {
	log           logging.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewCreateGroupHandler(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *createGroupHandler {
	return &createGroupHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (c *createGroupHandler) Handle(ctx context.Context, command *CreateGroupCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createGroupHandler.Handle")
	defer span.Finish()
	createDTO := &kafkaMessages.GroupCreate{
		ID:          command.CreateDto.ID.String(),
		Name:        command.CreateDto.Name,
		Description: command.CreateDto.Description,
		CreatorID:   command.CreateDto.CreatorID.String(),
		Active:      true,
	}
	dtoBytes, err := proto.Marshal(createDTO)
	if err != nil {
		return err
	}
	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.GroupCreate.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}

// UpdateGroupCmdHandler ...
type UpdateGroupCmdHandler interface {
	Handle(ctx context.Context, command *UpdateGroupCommand) error
}

type updateGroupCmdHandler struct {
	log           logging.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewUpdateGroupHandler(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *updateGroupCmdHandler {
	return &updateGroupCmdHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (c *updateGroupCmdHandler) Handle(ctx context.Context, command *UpdateGroupCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateGroupCmdHandler.Handle")
	defer span.Finish()
	updateDTO := &kafkaMessages.GroupUpdate{
		ID:          command.UpdateDto.ID.String(),
		Name:        command.UpdateDto.Name,
		Description: command.UpdateDto.Description,
	}
	dtoBytes, err := proto.Marshal(updateDTO)
	if err != nil {
		return err
	}
	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.GroupUpdate.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}

// DeleteGroupCmdHandler ...
type DeleteGroupCmdHandler interface {
	Handle(ctx context.Context, command *DeleteGroupCommand) error
}

type deleteGroupHandler struct {
	log           logging.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewDeleteGroupHandler(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *deleteGroupHandler {
	return &deleteGroupHandler{log: log, cfg: cfg, kafkaProducer: kafkaProducer}
}

func (c *deleteGroupHandler) Handle(ctx context.Context, command *DeleteGroupCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteGroupHandler.Handle")
	defer span.Finish()
	deleteDTO := &kafkaMessages.GroupDelete{ID: command.ID.String()}
	dtoBytes, err := proto.Marshal(deleteDTO)
	if err != nil {
		return err
	}
	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.GroupDelete.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}
