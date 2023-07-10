package queries

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/models"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	"github.com/JECSand/identity-service/pkg/logging"
)

/*
GET GROUP BY ID
*/

// GetGroupByIdHandler ...
type GetGroupByIdHandler interface {
	Handle(ctx context.Context, query *GetGroupByIdQuery) (*models.Group, error)
}

type getGroupByIdHandler struct {
	log    logging.Logger
	cfg    *config.Config
	pgRepo repositories.Repository
}

// NewGetGroupByIdHandler ...
func NewGetGroupByIdHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository) *getGroupByIdHandler {
	return &getGroupByIdHandler{
		log:    log,
		cfg:    cfg,
		pgRepo: pgRepo,
	}
}

// Handle ...
func (q *getGroupByIdHandler) Handle(ctx context.Context, query *GetGroupByIdQuery) (*models.Group, error) {
	return q.pgRepo.GetGroupById(ctx, query.ID)
}

/*
COUNT GROUPS
*/

// CountGroupsHandler ...
type CountGroupsHandler interface {
	Handle(ctx context.Context) (int, error)
}

type countGroupsHandler struct {
	log    logging.Logger
	cfg    *config.Config
	pgRepo repositories.Repository
}

// NewCountGroupsHandler ...
func NewCountGroupsHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository) *countGroupsHandler {
	return &countGroupsHandler{
		log:    log,
		cfg:    cfg,
		pgRepo: pgRepo,
	}
}

// Handle ...
func (q *countGroupsHandler) Handle(ctx context.Context) (int, error) {
	return q.pgRepo.CountGroups(ctx)
}
