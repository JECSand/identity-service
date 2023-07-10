package kafka

import (
	"context"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/tracing"
	kafkaMessages "github.com/JECSand/identity-service/protos/kafka"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/events"
	"github.com/JECSand/identity-service/query_service/identity/metrics"
	"github.com/JECSand/identity-service/query_service/identity/services"
	"github.com/avast/retry-go"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

const (
	PoolSize = 30
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

type queryMessageProcessor struct {
	log     logging.Logger
	cfg     *config.Config
	v       *validator.Validate
	us      *services.UserService
	gs      *services.GroupService
	ms      *services.MembershipService
	as      *services.AuthService
	metrics *metrics.QueryServiceMetrics
}

func NewQueryMessageProcessor(
	log logging.Logger,
	cfg *config.Config,
	v *validator.Validate,
	us *services.UserService,
	gs *services.GroupService,
	ms *services.MembershipService,
	as *services.AuthService,
	metrics *metrics.QueryServiceMetrics,
) *queryMessageProcessor {
	return &queryMessageProcessor{
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

func (s *queryMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
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
		case s.cfg.KafkaTopics.UserCreated.TopicName:
			s.processUserCreated(ctx, r, m)
		case s.cfg.KafkaTopics.UserUpdated.TopicName:
			s.processUserUpdated(ctx, r, m)
		case s.cfg.KafkaTopics.UserDeleted.TopicName:
			s.processUserDeleted(ctx, r, m)
		case s.cfg.KafkaTopics.GroupCreated.TopicName:
			s.processGroupCreated(ctx, r, m)
		case s.cfg.KafkaTopics.GroupUpdated.TopicName:
			s.processGroupUpdated(ctx, r, m)
		case s.cfg.KafkaTopics.GroupDeleted.TopicName:
			s.processGroupDeleted(ctx, r, m)
		case s.cfg.KafkaTopics.MembershipCreated.TopicName:
			s.processMembershipCreated(ctx, r, m)
		case s.cfg.KafkaTopics.MembershipUpdated.TopicName:
			s.processMembershipUpdated(ctx, r, m)
		case s.cfg.KafkaTopics.MembershipDeleted.TopicName:
			s.processMembershipDeleted(ctx, r, m)
		case s.cfg.KafkaTopics.TokenBlacklisted.TopicName:
			s.processBlacklistedToken(ctx, r, m)
		case s.cfg.KafkaTopics.PasswordUpdated.TopicName:
			s.processPasswordUpdated(ctx, r, m)
		}
	}
}

func (s *queryMessageProcessor) processMembershipCreated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.CreateMembershipKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processMembershipCreated")
	defer span.Finish()
	msg := &kafkaMessages.MembershipCreated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	createdMembership := events.NewCreatedMembership(
		msg.GetMembership().GetID(),
		msg.GetMembership().GetUserID(),
		msg.GetMembership().GetGroupID(),
		enums.MembershipStatus(msg.GetMembership().GetStatus()),
		enums.Role(msg.GetMembership().GetRole()),
		msg.GroupMembership.GetCreatedAt().AsTime(),
		msg.GroupMembership.GetUpdatedAt().AsTime(),
	)
	createdUserMembership := events.NewCreatedUserMembership(
		msg.GetUserMembership().GetID(),
		msg.GetUserMembership().GetGroupID(),
		msg.GetUserMembership().GetUserID(),
		msg.GetUserMembership().GetMembershipID(),
		msg.GetUserMembership().GetEmail(),
		msg.GetUserMembership().GetUsername(),
		enums.MembershipStatus(msg.GetUserMembership().GetStatus()),
		enums.Role(msg.GetUserMembership().GetRole()),
		msg.GroupMembership.GetCreatedAt().AsTime(),
		msg.GroupMembership.GetUpdatedAt().AsTime(),
	)
	createdGroupMembership := events.NewCreatedGroupMembership(
		msg.GetGroupMembership().GetID(),
		msg.GetGroupMembership().GetGroupID(),
		msg.GetGroupMembership().GetUserID(),
		msg.GetGroupMembership().GetMembershipID(),
		msg.GetGroupMembership().GetName(),
		msg.GetGroupMembership().GetDescription(),
		enums.MembershipStatus(msg.GetGroupMembership().GetStatus()),
		enums.Role(msg.GetGroupMembership().GetRole()),
		msg.GetGroupMembership().GetCreator(),
		msg.GroupMembership.GetCreatedAt().AsTime(),
		msg.GroupMembership.GetUpdatedAt().AsTime(),
	)
	event := events.NewCreateMembershipEvent(
		createdMembership,
		createdUserMembership,
		createdGroupMembership,
	)
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err := retry.Do(func() error {
		return s.ms.Events.CreateMembership.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("CreateMembership.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processMembershipUpdated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateMembershipKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processMembershipUpdated")
	defer span.Finish()
	msg := &kafkaMessages.MembershipUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	p := msg.GetMembership()
	event := events.NewUpdateMembershipEvent(p.GetID(), enums.MembershipStatus(p.GetStatus()), enums.Role(p.GetRole()), p.GetUpdatedAt().AsTime())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err := retry.Do(func() error {
		return s.ms.Events.UpdateMembership.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateMembership.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processMembershipDeleted(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.DeleteGroupKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processMembershipDeleted")
	defer span.Finish()
	msg := &kafkaMessages.MembershipDeleted{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	event := events.NewDeleteMembershipEvent(id)
	if err = s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.ms.Events.DeleteMembership.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("DeleteMembership.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processGroupCreated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.CreateGroupKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processGroupCreated")
	defer span.Finish()
	msg := &kafkaMessages.GroupCreated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	p := msg.GetGroup()
	event := events.NewCreateGroupEvent(
		p.GetID(),
		p.GetName(),
		p.GetDescription(),
		p.GetCreatorID(),
		p.GetActive(),
		p.GetCreatedAt().AsTime(),
		p.GetUpdatedAt().AsTime(),
	)
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err := retry.Do(func() error {
		return s.gs.Events.CreateGroup.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("CreateGroup.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processGroupUpdated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateGroupKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processGroupUpdated")
	defer span.Finish()
	msg := &kafkaMessages.GroupUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	p := msg.GetGroup()
	event := events.NewUpdateGroupEvent(p.GetID(), p.GetName(), p.GetDescription(), p.GetUpdatedAt().AsTime())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err := retry.Do(func() error {
		return s.gs.Events.UpdateGroup.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateGroup.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processGroupDeleted(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.DeleteGroupKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processGroupDeleted")
	defer span.Finish()
	msg := &kafkaMessages.GroupDeleted{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	event := events.NewDeleteGroupEvent(id)
	if err = s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.gs.Events.DeleteGroup.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("DeleteGroup.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processBlacklistedToken(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.BlacklistTokenKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processBlacklistedToken")
	defer span.Finish()
	msg := &kafkaMessages.TokenBlacklisted{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	p := msg.GetBlacklist()
	event := events.NewBlacklistTokenEvent(p.GetID(), p.GetAccessToken(), p.GetCreatedAt().AsTime(), p.GetUpdatedAt().AsTime())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err := retry.Do(func() error {
		return s.as.Events.BlacklistToken.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("BlacklistToken.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processPasswordUpdated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdatePasswordKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processPasswordUpdated")
	defer span.Finish()
	msg := &kafkaMessages.PasswordUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	event := events.NewUpdatePasswordEvent(msg.GetID(), msg.NewPassword, msg.GetUpdatedAt().AsTime())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err := retry.Do(func() error {
		return s.as.Events.UpdatePassword.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdatePassword.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processUserCreated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.CreateUserKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processUserCreated")
	defer span.Finish()
	msg := &kafkaMessages.UserCreated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	p := msg.GetUser()
	// TODO: Write logic for Root and Active User fields below
	event := events.NewCreateUserEvent(p.GetID(), p.GetEmail(), p.GetUsername(), p.GetPassword(), p.GetRoot(), p.GetActive(), p.GetCreatedAt().AsTime(), p.GetUpdatedAt().AsTime())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err := retry.Do(func() error {
		return s.us.Events.CreateUser.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("CreateUser.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processUserUpdated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateUserKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processUserUpdated")
	defer span.Finish()
	msg := &kafkaMessages.UserUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	p := msg.GetUser()
	event := events.NewUpdateUserEvent(p.GetID(), p.GetEmail(), p.GetUsername(), p.GetUpdatedAt().AsTime())
	if err := s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err := retry.Do(func() error {
		return s.us.Events.UpdateUser.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateUser.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) processUserDeleted(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.DeleteUserKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "queryMessageProcessor.processUserDeleted")
	defer span.Finish()
	msg := &kafkaMessages.UserDeleted{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	id, err := uuid.FromString(msg.GetID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	event := events.NewDeleteUserEvent(id)
	if err = s.v.StructCtx(ctx, event); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	if err = retry.Do(func() error {
		return s.us.Events.DeleteUser.Handle(ctx, event)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("DeleteUser.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}
	s.commitMessage(ctx, r, m)
}

func (s *queryMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.SuccessKafkaMessages.Inc()
	s.log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}

func (s *queryMessageProcessor) logProcessMessage(m kafka.Message, workerID int) {
	s.log.KafkaProcessMessage(m.Topic, m.Partition, string(m.Value), workerID, m.Offset, m.Time)
}

func (s *queryMessageProcessor) commitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.ErrorKafkaMessages.Inc()
	s.log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}
