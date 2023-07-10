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

// CreateMembershipCmdHandler ...
type CreateMembershipCmdHandler interface {
	Handle(ctx context.Context, command *CreateMembershipCommand) error
}

type createMembershipHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewCreateMembershipHandler ...
func NewCreateMembershipHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *createMembershipHandler {
	return &createMembershipHandler{
		log:           log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *createMembershipHandler) Handle(ctx context.Context, command *CreateMembershipCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createMembershipHandler.Handle")
	defer span.Finish()
	membershipDTO := &models.Membership{
		ID:      command.ID,
		UserID:  command.UserID,
		GroupID: command.GroupID,
		Status:  command.Status,
		Role:    command.Role,
	}
	membership, err := c.pgRepo.CreateMembership(ctx, membershipDTO)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)
	umChan := make(chan *models.UserMembership)
	gmChan := make(chan *models.GroupMembership)
	defer func() {
		cancel()
		close(errChan)
		close(umChan)
		close(gmChan)
	}()
	go func() {
		userMembership, err := c.pgRepo.GetUserMembershipById(ctx, membership.ID)
		select {
		case <-ctx.Done():
			return
		default:
		}
		umChan <- userMembership
		errChan <- err
	}()
	go func() {
		groupMembership, err := c.pgRepo.GetGroupMembershipById(ctx, membership.ID)
		select {
		case <-ctx.Done():
			return
		default:
		}
		gmChan <- groupMembership
		errChan <- err
	}()
	var userMembership *models.UserMembership
	var groupMembership *models.GroupMembership
	for i := 0; i < 4; i++ {
		select {
		case um := <-umChan:
			userMembership = um
		case gm := <-gmChan:
			groupMembership = gm
		case err = <-errChan:
			if err != nil {
				return err
			}
		}
	}
	msg := &kafkaMessages.MembershipCreated{
		Membership:      mappings.MembershipToGrpcMessage(membership),
		UserMembership:  mappings.UserMembershipToGrpcMessage(userMembership),
		GroupMembership: mappings.GroupMembershipToGrpcMessage(groupMembership),
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.MembershipCreated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}

// UpdateMembershipCmdHandler ...
type UpdateMembershipCmdHandler interface {
	Handle(ctx context.Context, command *UpdateMembershipCommand) error
}

type updateMembershipHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewUpdateMembershipHandler ...
func NewUpdateMembershipHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *updateMembershipHandler {
	return &updateMembershipHandler{log: log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *updateMembershipHandler) Handle(ctx context.Context, command *UpdateMembershipCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateMembershipHandler.Handle")
	defer span.Finish()
	membershipDTO := &models.Membership{
		ID:     command.ID,
		Status: command.Status,
		Role:   command.Role,
	}
	user, err := c.pgRepo.UpdateMembership(ctx, membershipDTO)
	if err != nil {
		return err
	}
	msg := &kafkaMessages.MembershipUpdated{Membership: mappings.MembershipToGrpcMessage(user)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.MembershipUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}

// DeleteMembershipCmdHandler ...
type DeleteMembershipCmdHandler interface {
	Handle(ctx context.Context, command *DeleteMembershipCommand) error
}

type deleteMembershipHandler struct {
	log           logging.Logger
	cfg           *config.Config
	pgRepo        repositories.Repository
	kafkaProducer kafkaClient.Producer
}

// NewDeleteMembershipHandler ...
func NewDeleteMembershipHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *deleteMembershipHandler {
	return &deleteMembershipHandler{
		log:           log,
		cfg:           cfg,
		pgRepo:        pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

// Handle ...
func (c *deleteMembershipHandler) Handle(ctx context.Context, command *DeleteMembershipCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteMembershipHandler.Handle")
	defer span.Finish()
	if err := c.pgRepo.DeleteMembershipById(ctx, command.ID); err != nil {
		return err
	}
	msg := &kafkaMessages.MembershipDeleted{ID: command.ID.String()}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.MembershipDeleted.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}
	return c.kafkaProducer.PublishMessage(ctx, message)
}
