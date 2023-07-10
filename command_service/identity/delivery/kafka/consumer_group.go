package kafka

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	"github.com/JECSand/identity-service/command_service/identity/metrics"
	"github.com/JECSand/identity-service/command_service/identity/services"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	kafkaMessages "github.com/JECSand/identity-service/protos/kafka"
	"github.com/avast/retry-go"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

const (
	PoolSize      = 30
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

type identityMessageProcessor struct {
	log     logging.Logger
	cfg     *config.Config
	v       *validator.Validate
	us      *services.UserService
	gs      *services.GroupService
	ms      *services.MembershipService
	as      *services.AuthService
	metrics *metrics.CommandServiceMetrics
}

func NewIdentityMessageProcessor(
	log logging.Logger,
	cfg *config.Config,
	v *validator.Validate,
	us *services.UserService,
	gs *services.GroupService,
	ms *services.MembershipService,
	as *services.AuthService,
	metrics *metrics.CommandServiceMetrics,
) *identityMessageProcessor {
	return &identityMessageProcessor{
		log:     log,
		cfg:     cfg,
		v:       v,
		us:      us,
		gs:      gs,
		ms:      ms,
		as:      as,
		metrics: metrics,
	}
}

func (s *identityMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.SuccessKafkaMessages.Inc()
	s.log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}

func (s *identityMessageProcessor) commitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.ErrorKafkaMessages.Inc()
	s.log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}

func (s *identityMessageProcessor) logProcessMessage(m kafka.Message, workerID int) {
	s.log.KafkaProcessMessage(m.Topic, m.Partition, string(m.Value), workerID, m.Offset, m.Time)
}

func (s *identityMessageProcessor) processBlacklistToken(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.BlacklistTokenKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processBlacklistToken")
	defer span.Finish()
	var msg kafkaMessages.TokenBlacklist
	if err := proto.Unmarshal(m.Value, &msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewBlacklistTokenCommand(id, msg.GetAccessToken())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.as.Commands.BlacklistToken.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("BlacklistToken.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processUpdatePassword(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.PasswordUpdateKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processUpdatePassword")
	defer span.Finish()
	msg := &kafkaMessages.PasswordUpdate{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewUpdatePasswordCommand(id, msg.GetCurrentPassword(), msg.GetNewPassword())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.as.Commands.UpdatePassword.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdatePassword.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processCreateGroup(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.CreateGroupKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processCreateGroup")
	defer span.Finish()
	var msg kafkaMessages.GroupCreate
	if err := proto.Unmarshal(m.Value, &msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	creatorId, err := uuid.FromString(msg.GetCreatorID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewCreateGroupCommand(id, msg.GetName(), msg.GetDescription(), creatorId, msg.GetActive())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.gs.Commands.CreateGroup.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("CreateGroup.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processUpdateGroup(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateGroupKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processUpdateGroup")
	defer span.Finish()
	msg := &kafkaMessages.GroupUpdate{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewUpdateGroupCommand(id, msg.GetName(), msg.GetDescription())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.gs.Commands.UpdateGroup.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateGroup.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processDeleteGroup(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.DeleteGroupKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processDeleteGroup")
	defer span.Finish()
	msg := &kafkaMessages.GroupDelete{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewDeleteGroupCommand(id)
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.gs.Commands.DeleteGroup.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("DeleteGroup.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processCreateMembership(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.CreateMembershipKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processCreateMembership")
	defer span.Finish()
	var msg kafkaMessages.MembershipCreate
	if err := proto.Unmarshal(m.Value, &msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	userId, err := uuid.FromString(msg.GetUserID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	groupId, err := uuid.FromString(msg.GetGroupID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewCreateMembershipCommand(id, userId, groupId, enums.MembershipStatus(msg.GetStatus()), enums.Role(msg.GetRole()))
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.ms.Commands.CreateMembership.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("CreateMembership.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processUpdateMembership(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateMembershipKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processUpdateMembership")
	defer span.Finish()
	msg := &kafkaMessages.MembershipUpdate{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewUpdateMembershipCommand(id, enums.MembershipStatus(msg.GetStatus()), enums.Role(msg.GetRole()))
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.ms.Commands.UpdateMembership.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateMembership.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processDeleteMembership(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.DeleteGroupKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processDeleteMembership")
	defer span.Finish()
	msg := &kafkaMessages.MembershipDelete{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewDeleteMembershipCommand(id)
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.ms.Commands.DeleteMembership.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("DeleteMembership.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processCreateUser(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.CreateUserKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processCreateUser")
	defer span.Finish()
	var msg kafkaMessages.UserCreate
	if err := proto.Unmarshal(m.Value, &msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewCreateUserCommand(id, msg.GetEmail(), msg.GetUsername(), msg.GetPassword(), false, msg.GetActive())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.us.Commands.CreateUser.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("CreateUser.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processUpdateUser(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateUserKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processUpdateUser")
	defer span.Finish()
	msg := &kafkaMessages.UserUpdate{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewUpdateUserCommand(id, msg.GetEmail(), msg.GetUsername())
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.us.Commands.UpdateUser.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateUser.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) processDeleteUser(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.DeleteUserKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "identityMessageProcessor.processDeleteUser")
	defer span.Finish()
	msg := &kafkaMessages.UserDelete{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	command := commands.NewDeleteUserCommand(id)
	if err = s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.us.Commands.DeleteUser.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("DeleteUser.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *identityMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		m, err := r.FetchMessage(ctx)
		if err != nil {
			s.log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}
		s.logProcessMessage(m, workerID)
		switch m.Topic {
		case s.cfg.KafkaTopics.UserCreate.TopicName:
			s.processCreateUser(ctx, r, m)
		case s.cfg.KafkaTopics.UserUpdate.TopicName:
			s.processUpdateUser(ctx, r, m)
		case s.cfg.KafkaTopics.UserDelete.TopicName:
			s.processDeleteUser(ctx, r, m)
		case s.cfg.KafkaTopics.TokenBlacklist.TopicName:
			s.processBlacklistToken(ctx, r, m)
		case s.cfg.KafkaTopics.PasswordUpdate.TopicName:
			s.processUpdatePassword(ctx, r, m)
		}
	}
}
