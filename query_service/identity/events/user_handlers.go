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

// CreateUserEventHandler ...
type CreateUserEventHandler interface {
	Handle(ctx context.Context, event *CreateUserEvent) error
}

type createUserEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewCreateUserEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *createUserEventHandler {
	return &createUserEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *createUserEventHandler) Handle(ctx context.Context, event *CreateUserEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createUserEventHandler.Handle")
	defer span.Finish()
	user := &entities.User{
		ID:        event.ID,
		Email:     event.Email,
		Username:  event.Username,
		Password:  event.Password,
		Root:      event.Root,
		Active:    event.Active,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
	created, err := c.mongoDB.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	c.redisCache.PutUser(ctx, created.ID, created)
	return nil
}

// UpdateUserEventHandler ...
type UpdateUserEventHandler interface {
	Handle(ctx context.Context, event *UpdateUserEvent) error
}

type updateUserEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewUpdateUserEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *updateUserEventHandler {
	return &updateUserEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *updateUserEventHandler) Handle(ctx context.Context, event *UpdateUserEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateUserEventHandler.Handle")
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)
	defer func() {
		span.Finish()
		cancel()
		close(errChan)
	}()
	user := &entities.User{
		ID:        event.ID,
		Email:     event.Email,
		Username:  event.Username,
		UpdatedAt: event.UpdatedAt,
	}
	go func() {
		updated, err := c.mongoDB.UpdateUser(ctx, user)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			c.redisCache.PutUser(ctx, updated.ID, updated)
		}
		errChan <- err
	}()
	go func() {
		filter := &entities.UserMembership{
			UserID: user.ID,
		}
		update := &entities.UserMembership{
			Username:  user.Username,
			Email:     user.Email,
			UpdatedAt: user.UpdatedAt,
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

// DeleteUserEventHandler ...
type DeleteUserEventHandler interface {
	Handle(ctx context.Context, event *DeleteUserEvent) error
}

type deleteUserEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewDeleteUserEventHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *deleteUserEventHandler {
	return &deleteUserEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *deleteUserEventHandler) Handle(ctx context.Context, event *DeleteUserEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteUserEventHandler.Handle")
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)
	defer func() {
		span.Finish()
		cancel()
		close(errChan)
	}()
	go func() {
		err := c.mongoDB.DeleteUser(ctx, event.ID)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			c.redisCache.DeleteUser(ctx, event.ID.String())
		}
		errChan <- err
	}()
	go func() {
		filter := &entities.Membership{
			UserID: event.ID.String(),
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
			UserID: event.ID.String(),
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
			UserID: event.ID.String(),
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
	return nil
}
