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

// GetMembershipByIdHandler ...
type GetMembershipByIdHandler interface {
	Handle(ctx context.Context, query *GetMembershipByIdQuery) (*entities.Membership, error)
}

type getMembershipByIdHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewGetMembershipByIdHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *getMembershipByIdHandler {
	return &getMembershipByIdHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (q *getMembershipByIdHandler) Handle(ctx context.Context, query *GetMembershipByIdQuery) (*entities.Membership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getMembershipByIdHandler.Handle")
	defer span.Finish()
	if membership, err := q.redisCache.GetMembership(ctx, query.ID.String()); err == nil && membership != nil {
		return membership, nil
	}
	membership, err := q.mongoDB.GetMembershipById(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	q.redisCache.PutMembership(ctx, membership.ID, membership)
	return membership, nil
}

// GetGroupMembershipHandler ...
type GetGroupMembershipHandler interface {
	Handle(ctx context.Context, query *GetGroupMembershipQuery) (*entities.GroupMembershipsList, error)
}

type getGroupMembershipHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewGetGroupMembershipHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *getGroupMembershipHandler {
	return &getGroupMembershipHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (s *getGroupMembershipHandler) Handle(ctx context.Context, query *GetGroupMembershipQuery) (*entities.GroupMembershipsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getGroupMembershipHandler.Handle")
	defer span.Finish()
	return s.mongoDB.GetGroupMembershipByUserId(ctx, query.UserID, query.Pagination)
}

// GetUserMembershipHandler ...
type GetUserMembershipHandler interface {
	Handle(ctx context.Context, query *GetUserMembershipQuery) (*entities.UserMembershipsList, error)
}

type getUserMembershipHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewGetUserMembershipHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *getUserMembershipHandler {
	return &getUserMembershipHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (s *getUserMembershipHandler) Handle(ctx context.Context, query *GetUserMembershipQuery) (*entities.UserMembershipsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getUserMembershipHandler.Handle")
	defer span.Finish()
	return s.mongoDB.GetUserMembershipByGroupId(ctx, query.GroupID, query.Pagination)
}
