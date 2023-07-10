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
	"time"
)

// CreateGroupCmdHandler ...
type CreateGroupCmdHandler interface {
	Handle(ctx context.Context, command *CreateGroupCommand) error
}

type createGroupHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewCreateGroupHandler ...
func NewCreateGroupHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *createGroupHandler {
	return &createGroupHandler{
		log:           log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *createGroupHandler) Handle(ctx context.Context, command *CreateGroupCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createGroupHandler.Handle")
	defer span.Finish()
	groupDTO := &models.Group{
		ID:          command.ID,
		Name:        command.Name,
		Description: command.Description,
		CreatorID:   command.CreatorID,
		Active:      command.Active,
	}
	group, err := c.pgRepo.CreateGroup(ctx, groupDTO)
	if err != nil {
		return err
	}
	msg := &kafkaMessages.GroupCreated{Group: mappings.GroupToGrpcMessage(group)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.GroupCreated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}

// UpdateGroupCmdHandler ...
type UpdateGroupCmdHandler interface {
	Handle(ctx context.Context, command *UpdateGroupCommand) error
}

type updateGroupHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewUpdateGroupHandler ...
func NewUpdateGroupHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *updateGroupHandler {
	return &updateGroupHandler{log: log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *updateGroupHandler) Handle(ctx context.Context, command *UpdateGroupCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateGroupHandler.Handle")
	defer span.Finish()
	groupDTO := &models.Group{
		ID:          command.ID,
		Name:        command.Name,
		Description: command.Description,
	}
	user, err := c.pgRepo.UpdateGroup(ctx, groupDTO)
	if err != nil {
		return err
	}
	msg := &kafkaMessages.GroupUpdated{Group: mappings.GroupToGrpcMessage(user)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.GroupUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}

// DeleteGroupCmdHandler ...
type DeleteGroupCmdHandler interface {
	Handle(ctx context.Context, command *DeleteGroupCommand) error
}

type deleteGroupHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewDeleteGroupHandler ...
func NewDeleteGroupHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *deleteGroupHandler {
	return &deleteGroupHandler{
		log:           log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *deleteGroupHandler) Handle(ctx context.Context, command *DeleteGroupCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteGroupHandler.Handle")
	defer span.Finish()
	if err := c.pgRepo.DeleteGroupById(ctx, command.ID); err != nil {
		return err
	}
	msg := &kafkaMessages.GroupDeleted{ID: command.ID.String()}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.GroupDeleted.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}
