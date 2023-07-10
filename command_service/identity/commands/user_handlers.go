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

// CreateUserCmdHandler ...
type CreateUserCmdHandler interface {
	Handle(ctx context.Context, command *CreateUserCommand) error
}

type createUserHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewCreateUserHandler ...
func NewCreateUserHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *createUserHandler {
	return &createUserHandler{
		log:           log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *createUserHandler) Handle(ctx context.Context, command *CreateUserCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createUserHandler.Handle")
	defer span.Finish()
	userDTO := &models.User{
		ID:       command.ID,
		Email:    command.Email,
		Username: command.Username,
		Password: command.Password,
		Root:     command.Root,
		Active:   command.Active,
	}
	if err := userDTO.HashPassword(); err != nil {
		return err
	}
	user, err := c.pgRepo.CreateUser(ctx, userDTO)
	if err != nil {
		return err
	}
	msg := &kafkaMessages.UserCreated{User: mappings.UserToGrpcMessage(user)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.UserCreated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}

// UpdateUserCmdHandler ...
type UpdateUserCmdHandler interface {
	Handle(ctx context.Context, command *UpdateUserCommand) error
}

type updateUserHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewUpdateUserHandler ...
func NewUpdateUserHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *updateUserHandler {
	return &updateUserHandler{log: log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *updateUserHandler) Handle(ctx context.Context, command *UpdateUserCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateUserHandler.Handle")
	defer span.Finish()
	userDTO := &models.User{
		ID:       command.ID,
		Email:    command.Email,
		Username: command.Username,
	}
	user, err := c.pgRepo.UpdateUser(ctx, userDTO)
	if err != nil {
		return err
	}
	msg := &kafkaMessages.UserUpdated{User: mappings.UserToGrpcMessage(user)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.UserUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}

// DeleteUserCmdHandler ...
type DeleteUserCmdHandler interface {
	Handle(ctx context.Context, command *DeleteUserCommand) error
}

type deleteUserHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewDeleteUserHandler ...
func NewDeleteUserHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *deleteUserHandler {
	return &deleteUserHandler{
		log:           log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *deleteUserHandler) Handle(ctx context.Context, command *DeleteUserCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteUserHandler.Handle")
	defer span.Finish()
	if err := c.pgRepo.DeleteUserById(ctx, command.ID); err != nil {
		return err
	}
	msg := &kafkaMessages.UserDeleted{ID: command.ID.String()}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.UserDeleted.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}
