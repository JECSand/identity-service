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

// BlacklistTokenEventHandler ...
type BlacklistTokenEventHandler interface {
	Handle(ctx context.Context, event *BlacklistTokenEvent) error
}

type blacklistTokenEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewBlacklistTokenEventHandler(
	log logging.Logger,
	cfg *config.Config,
	mongoDB data.Database,
	redisCache cache.Cache,
) *blacklistTokenEventHandler {
	return &blacklistTokenEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *blacklistTokenEventHandler) Handle(ctx context.Context, event *BlacklistTokenEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blacklistTokenEventHandler.Handle")
	defer span.Finish()
	user := &entities.Blacklist{
		ID:          event.ID,
		AccessToken: event.AccessToken,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
	created, err := c.mongoDB.BlacklistToken(ctx, user)
	if err != nil {
		return err
	}
	c.redisCache.PutToken(ctx, created.ID, created)
	return nil
}

// UpdatePasswordEventHandler ...
type UpdatePasswordEventHandler interface {
	Handle(ctx context.Context, event *UpdatePasswordEvent) error
}

type updatePasswordEventHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewUpdatePasswordEventHandler(
	log logging.Logger,
	cfg *config.Config,
	mongoDB data.Database,
	redisCache cache.Cache,
) *updatePasswordEventHandler {
	return &updatePasswordEventHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (c *updatePasswordEventHandler) Handle(ctx context.Context, event *UpdatePasswordEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updatePasswordEventHandler.Handle")
	defer span.Finish()
	user := &entities.User{
		ID:        event.ID,
		Password:  event.NewPassword,
		UpdatedAt: event.UpdatedAt,
	}
	updated, err := c.mongoDB.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	c.redisCache.PutUser(ctx, updated.ID, updated)
	return nil
}
