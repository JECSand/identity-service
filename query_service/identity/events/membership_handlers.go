package events

import (
	"context"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/cache"
	"github.com/JECSand/identity-service/query_service/identity/data"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/opentracing/opentracing-go"
)

// CreateMembershipEventHandler ...
type CreateMembershipEventHandler interface {
	Handle(ctx context.Context, event *CreateMembershipEvent) error
}

type createMembershipEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewCreateMembershipEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *createMembershipEventHandler {
	return &createMembershipEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *createMembershipEventHandler) Handle(ctx context.Context, event *CreateMembershipEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createMembershipEventHandler.Handle")
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)
	defer func() {
		span.Finish()
		cancel()
		close(errChan)
	}()
	go func() {
		membership := &entities.Membership{
			ID:        event.Membership.ID,
			UserID:    event.Membership.UserID,
			GroupID:   event.Membership.GroupID,
			Status:    event.Membership.Status,
			Role:      event.Membership.Role,
			CreatedAt: event.Membership.CreatedAt,
			UpdatedAt: event.Membership.UpdatedAt,
		}
		createdMembership, err := c.mongoDB.CreateMembership(ctx, membership)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			c.redisCache.PutMembership(ctx, createdMembership.ID, createdMembership)
		}
		errChan <- err
	}()
	go func() {
		userMembership := &entities.UserMembership{
			ID:           event.UserMembership.ID,
			GroupID:      event.UserMembership.GroupID,
			UserID:       event.UserMembership.UserID,
			MembershipID: event.UserMembership.MembershipID,
			Email:        event.UserMembership.Email,
			Username:     event.UserMembership.Username,
			Status:       event.UserMembership.Status,
			Role:         event.UserMembership.Role,
			CreatedAt:    event.UserMembership.CreatedAt,
			UpdatedAt:    event.UserMembership.UpdatedAt,
		}
		_, err := c.mongoDB.CreateUserMembership(ctx, userMembership)
		select {
		case <-ctx.Done():
			return
		default:
		}
		//if err == nil {
		//	c.redisCache.PutUserMembership(ctx, createdUserMembership.ID, createdUserMembership)
		//}
		errChan <- err
	}()
	go func() {
		groupMembership := &entities.GroupMembership{
			ID:           event.GroupMembership.ID,
			UserID:       event.GroupMembership.UserID,
			GroupID:      event.GroupMembership.GroupID,
			MembershipID: event.GroupMembership.MembershipID,
			Name:         event.GroupMembership.Name,
			Description:  event.GroupMembership.Description,
			Status:       event.GroupMembership.Status,
			Role:         event.GroupMembership.Role,
			Creator:      event.GroupMembership.Creator,
			CreatedAt:    event.GroupMembership.CreatedAt,
			UpdatedAt:    event.GroupMembership.UpdatedAt,
		}
		_, err := c.mongoDB.CreateGroupMembership(ctx, groupMembership)
		select {
		case <-ctx.Done():
			return
		default:
		}
		//if err == nil {
		//	c.redisCache.PutGroupMembership(ctx, createdGroupMembership.ID, createdGroupMembership)
		//}
		errChan <- err
	}()
	var err error
	for i := 0; i < 3; i++ {
		select {
		case err = <-errChan:
			if err != nil {
				break
			}
		}
	}
	return err
}

// UpdateMembershipEventHandler ...
type UpdateMembershipEventHandler interface {
	Handle(ctx context.Context, event *UpdateMembershipEvent) error
}

type updateMembershipEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewUpdateMembershipEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *updateMembershipEventHandler {
	return &updateMembershipEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *updateMembershipEventHandler) Handle(ctx context.Context, event *UpdateMembershipEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateMembershipEventHandler.Handle")
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)
	defer func() {
		span.Finish()
		cancel()
		close(errChan)
	}()
	go func() {
		membership := &entities.Membership{
			ID:        event.ID,
			Status:    event.Status,
			Role:      event.Role,
			UpdatedAt: event.UpdatedAt,
		}
		updated, err := c.mongoDB.UpdateMembership(ctx, membership)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			c.redisCache.PutMembership(ctx, updated.ID, updated)
		}
		errChan <- err
	}()
	go func() {
		filter := &entities.GroupMembership{
			MembershipID: event.ID,
		}
		update := &entities.GroupMembership{
			Status:    event.Status,
			Role:      event.Role,
			UpdatedAt: event.UpdatedAt,
		}
		err := c.mongoDB.UpdateGroupMemberships(ctx, filter, update)
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- err
	}()
	go func() {
		filter := &entities.UserMembership{
			MembershipID: event.ID,
		}
		update := &entities.UserMembership{
			Status:    event.Status,
			Role:      event.Role,
			UpdatedAt: event.UpdatedAt,
		}
		err := c.mongoDB.UpdateUserMemberships(ctx, filter, update)
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- err
	}()
	var err error
	for i := 0; i < 3; i++ {
		select {
		case err = <-errChan:
			if err != nil {
				break
			}
		}
	}
	return err
}

// DeleteMembershipEventHandler ...
type DeleteMembershipEventHandler interface {
	Handle(ctx context.Context, event *DeleteMembershipEvent) error
}

type deleteMembershipEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewDeleteMembershipEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *deleteMembershipEventHandler {
	return &deleteMembershipEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *deleteMembershipEventHandler) Handle(ctx context.Context, event *DeleteMembershipEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteMembershipEventHandler.Handle")
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)
	defer func() {
		span.Finish()
		cancel()
		close(errChan)
	}()
	go func() {
		err := c.mongoDB.DeleteMembership(ctx, event.ID)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			c.redisCache.DeleteMembership(ctx, event.ID.String())
		}
		errChan <- err
	}()
	go func() {
		err := c.mongoDB.DeleteUserMembershipByMembershipId(ctx, event.ID)
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- err
	}()
	go func() {
		err := c.mongoDB.DeleteGroupMembershipByMembershipId(ctx, event.ID)
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- err
	}()
	var err error
	for i := 0; i < 3; i++ {
		select {
		case err = <-errChan:
			if err != nil {
				break
			}
		}
	}
	return err
}
