package queries

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/models"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	"github.com/JECSand/identity-service/pkg/logging"
)

/*
CHECK TOKEN BLACKLIST
*/

// CheckTokenBlacklistHandler ...
type CheckTokenBlacklistHandler interface {
	Handle(ctx context.Context, query *CheckTokenBlacklistQuery) (*models.Blacklist, error)
}

type checkTokenBlacklistHandler struct {
	log    logging.Logger
	cfg    *config.Config
	pgRepo repositories.Repository
}

// NewCheckTokenBlacklistHandler ...
func NewCheckTokenBlacklistHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository) *checkTokenBlacklistHandler {
	return &checkTokenBlacklistHandler{
		log:    log,
		cfg:    cfg,
		pgRepo: pgRepo,
	}
}

// Handle ...
func (q *checkTokenBlacklistHandler) Handle(ctx context.Context, query *CheckTokenBlacklistQuery) (*models.Blacklist, error) {
	return q.pgRepo.CheckBlacklist(ctx, query.AccessToken)
}
