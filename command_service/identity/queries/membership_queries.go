package queries

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/models"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	"github.com/JECSand/identity-service/pkg/logging"
)

/*
GET MEMBERSHIP BY ID
*/

// GetMembershipByIdHandler ...
type GetMembershipByIdHandler interface {
	Handle(ctx context.Context, query *GetMembershipByIdQuery) (*models.Membership, error)
}

type getMembershipByIdHandler struct {
	log    logging.Logger
	cfg    *config.Config
	pgRepo repositories.Repository
}

// NewGetMembershipByIdHandler ...
func NewGetMembershipByIdHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository) *getMembershipByIdHandler {
	return &getMembershipByIdHandler{
		log:    log,
		cfg:    cfg,
		pgRepo: pgRepo,
	}
}

// Handle ...
func (q *getMembershipByIdHandler) Handle(ctx context.Context, query *GetMembershipByIdQuery) (*models.Membership, error) {
	return q.pgRepo.GetMembershipById(ctx, query.ID)
}

/*
GET USER MEMBERSHIP BY ID
*/

// GetUserMembershipByIdHandler ...
type GetUserMembershipByIdHandler interface {
	Handle(ctx context.Context, query *GetUserMembershipByIdQuery) (*models.UserMembership, error)
}

type getUserMembershipByIdHandler struct {
	log    logging.Logger
	cfg    *config.Config
	pgRepo repositories.Repository
}

// NewGetUserMembershipByIdHandler ...
func NewGetUserMembershipByIdHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository) *getUserMembershipByIdHandler {
	return &getUserMembershipByIdHandler{
		log:    log,
		cfg:    cfg,
		pgRepo: pgRepo,
	}
}

// Handle ...
func (q *getUserMembershipByIdHandler) Handle(ctx context.Context, query *GetUserMembershipByIdQuery) (*models.UserMembership, error) {
	return q.pgRepo.GetUserMembershipById(ctx, query.ID)
}

/*
GET GROUP MEMBERSHIP BY ID
*/

// GetGroupMembershipByIdHandler ...
type GetGroupMembershipByIdHandler interface {
	Handle(ctx context.Context, query *GetGroupMembershipByIdQuery) (*models.GroupMembership, error)
}

type getGroupMembershipByIdHandler struct {
	log    logging.Logger
	cfg    *config.Config
	pgRepo repositories.Repository
}

// NewGetGroupMembershipByIdHandler ...
func NewGetGroupMembershipByIdHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository) *getGroupMembershipByIdHandler {
	return &getGroupMembershipByIdHandler{
		log:    log,
		cfg:    cfg,
		pgRepo: pgRepo,
	}
}

// Handle ...
func (q *getGroupMembershipByIdHandler) Handle(ctx context.Context, query *GetGroupMembershipByIdQuery) (*models.GroupMembership, error) {
	return q.pgRepo.GetGroupMembershipById(ctx, query.ID)
}

/*
COUNT GROUPS
*/

// CountMembershipsHandler ...
type CountMembershipsHandler interface {
	Handle(ctx context.Context) (int, error)
}

type countMembershipsHandler struct {
	log    logging.Logger
	cfg    *config.Config
	pgRepo repositories.Repository
}

// NewCountMembershipsHandler ...
func NewCountMembershipsHandler(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository) *countMembershipsHandler {
	return &countMembershipsHandler{
		log:    log,
		cfg:    cfg,
		pgRepo: pgRepo,
	}
}

// Handle ...
func (q *countMembershipsHandler) Handle(ctx context.Context) (int, error) {
	return q.pgRepo.CountMemberships(ctx)
}
