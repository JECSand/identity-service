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

// CreateGroupEventHandler ...
type CreateGroupEventHandler interface {
	Handle(ctx context.Context, event *CreateGroupEvent) error
}

type createGroupEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewCreateGroupEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *createGroupEventHandler {
	return &createGroupEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *createGroupEventHandler) Handle(ctx context.Context, event *CreateGroupEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createGroupEventHandler.Handle")
	defer span.Finish()
	group := &entities.Group{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		CreatorID:   event.CreatorID,
		Active:      event.Active,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
	created, err := c.mongoDB.CreateGroup(ctx, group)
	if err != nil {
		return err
	}
	c.redisCache.PutGroup(ctx, created.ID, created)
	return nil
}

// UpdateGroupEventHandler ...
type UpdateGroupEventHandler interface {
	Handle(ctx context.Context, event *UpdateGroupEvent) error
}

type updateGroupEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewUpdateGroupEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *updateGroupEventHandler {
	return &updateGroupEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *updateGroupEventHandler) Handle(ctx context.Context, event *UpdateGroupEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateGroupEventHandler.Handle")
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)
	defer func() {
		span.Finish()
		cancel()
		close(errChan)
	}()
	group := &entities.Group{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		UpdatedAt:   event.UpdatedAt,
	}
	go func() {
		updated, err := c.mongoDB.UpdateGroup(ctx, group)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			c.redisCache.PutGroup(ctx, updated.ID, updated)
		}
		errChan <- err
	}()
	go func() {
		filter := &entities.GroupMembership{
			GroupID: group.ID,
		}
		update := &entities.GroupMembership{
			Name:        group.Name,
			Description: group.Description,
			UpdatedAt:   group.UpdatedAt,
		}
		err := c.mongoDB.UpdateGroupMemberships(ctx, filter, update)
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- err
	}()
	var err error
	for i := 0; i < 2; i++ {
		select {
		case err = <-errChan:
			if err != nil {
				break
			}
		}
	}
	return err
}

// DeleteGroupEventHandler ...
type DeleteGroupEventHandler interface {
	Handle(ctx context.Context, event *DeleteGroupEvent) error
}

type deleteGroupEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewDeleteGroupEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *deleteGroupEventHandler {
	return &deleteGroupEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *deleteGroupEventHandler) Handle(ctx context.Context, event *DeleteGroupEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteGroupEventHandler.Handle")
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)
	defer func() {
		span.Finish()
		cancel()
		close(errChan)
	}()
	go func() {
		err := c.mongoDB.DeleteGroup(ctx, event.ID)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			c.redisCache.DeleteGroup(ctx, event.ID.String())
		}
		errChan <- err
	}()
	go func() {
		filter := &entities.Membership{
			GroupID: event.ID.String(),
		}
		err := c.mongoDB.DeleteMemberships(ctx, filter)
		select {
		case <-ctx.Done():
			return
		default:
		}
		// TODO: Clear Cache of deleted Memberships
		//if err == nil {
		//	c.redisCache.DeleteMembership(ctx, event.ID.String())
		//}
		errChan <- err
	}()
	go func() {
		filter := &entities.GroupMembership{
			GroupID: event.ID.String(),
		}
		err := c.mongoDB.DeleteGroupMemberships(ctx, filter)
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- err
	}()
	go func() {
		filter := &entities.UserMembership{
			GroupID: event.ID.String(),
		}
		err := c.mongoDB.DeleteUserMemberships(ctx, filter)
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- err
	}()
	var err error
	for i := 0; i < 4; i++ {
		select {
		case err = <-errChan:
			if err != nil {
				break
			}
		}
	}
	return err
}
