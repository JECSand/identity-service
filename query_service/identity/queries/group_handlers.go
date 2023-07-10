package queries

import (
	"context"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/cache"
	"github.com/JECSand/identity-service/query_service/identity/data"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/opentracing/opentracing-go"
)

// GetGroupByIdHandler ...
type GetGroupByIdHandler interface {
	Handle(ctx context.Context, query *GetGroupByIdQuery) (*entities.Group, error)
}

type getGroupByIdHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewGetGroupByIdHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *getGroupByIdHandler {
	return &getGroupByIdHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (q *getGroupByIdHandler) Handle(ctx context.Context, query *GetGroupByIdQuery) (*entities.Group, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getGroupByIdHandler.Handle")
	defer span.Finish()
	if group, err := q.redisCache.GetGroup(ctx, query.ID.String()); err == nil && group != nil {
		return group, nil
	}
	group, err := q.mongoDB.GetGroupById(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	q.redisCache.PutGroup(ctx, group.ID, group)
	return group, nil
}

// SearchGroupHandler ...
type SearchGroupHandler interface {
	Handle(ctx context.Context, query *SearchGroupQuery) (*entities.GroupsList, error)
}

type searchGroupHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewSearchGroupHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *searchGroupHandler {
	return &searchGroupHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (s *searchGroupHandler) Handle(ctx context.Context, query *SearchGroupQuery) (*entities.GroupsList, error) {
	return s.mongoDB.SearchGroups(ctx, query.Text, query.Pagination)
}
