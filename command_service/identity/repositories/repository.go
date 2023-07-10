package repositories

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/models"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	blacklist   *blacklistRepository
	users       *userRepository
	groups      *groupRepository
	memberships *membershipRepository
}

// NewRepository ...
func NewRepository(log logging.Logger, cfg *config.Config, db *pgxpool.Pool) *repository {
	u := NewUserRepository(log, cfg, db)
	g := NewGroupRepository(log, cfg, db)
	m := NewMembershipRepository(log, cfg, db)
	b := NewBlacklistRepository(log, cfg, db)
	return &repository{
		blacklist:   b,
		users:       u,
		groups:      g,
		memberships: m,
	}
}

func (d *repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return d.users.Create(ctx, user)
}

func (d *repository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return d.users.Update(ctx, user)
}

func (d *repository) UpdateUserPassword(ctx context.Context, user *models.User) (*models.User, error) {
	return d.users.UpdatePassword(ctx, user)
}

func (d *repository) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	return d.users.DeleteByID(ctx, id)
}

func (d *repository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return d.users.GetById(ctx, id)
}

func (d *repository) CountUsers(ctx context.Context) (int, error) {
	return d.users.Count(ctx)
}

func (d *repository) CreateGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	return d.groups.Create(ctx, group)
}

func (d *repository) UpdateGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	return d.groups.Update(ctx, group)
}

func (d *repository) DeleteGroupById(ctx context.Context, id uuid.UUID) error {
	return d.groups.DeleteByID(ctx, id)
}

func (d *repository) GetGroupById(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	return d.groups.GetById(ctx, id)
}

func (d *repository) CountGroups(ctx context.Context) (int, error) {
	return d.groups.Count(ctx)
}

func (d *repository) CreateMembership(ctx context.Context, membership *models.Membership) (*models.Membership, error) {
	return d.memberships.Create(ctx, membership)
}

func (d *repository) UpdateMembership(ctx context.Context, membership *models.Membership) (*models.Membership, error) {
	return d.memberships.Update(ctx, membership)
}

func (d *repository) DeleteMembershipById(ctx context.Context, id uuid.UUID) error {
	return d.memberships.DeleteByID(ctx, id)
}

func (d *repository) GetMembershipById(ctx context.Context, id uuid.UUID) (*models.Membership, error) {
	return d.memberships.GetById(ctx, id)
}

func (d *repository) CountMemberships(ctx context.Context) (int, error) {
	return d.memberships.Count(ctx)
}

func (d *repository) GetUserMembershipById(ctx context.Context, id uuid.UUID) (*models.UserMembership, error) {
	return d.memberships.GetUserMembershipById(ctx, id)
}

func (d *repository) GetGroupMembershipById(ctx context.Context, id uuid.UUID) (*models.GroupMembership, error) {
	return d.memberships.GetGroupMembershipById(ctx, id)
}

func (d *repository) BlacklistToken(ctx context.Context, blacklist *models.Blacklist) (*models.Blacklist, error) {
	return d.blacklist.Create(ctx, blacklist)
}

func (d *repository) CheckBlacklist(ctx context.Context, accessToken string) (*models.Blacklist, error) {
	return d.blacklist.GetByAccessToken(ctx, accessToken)
}

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUserById(ctx context.Context, id uuid.UUID) error
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	CountUsers(ctx context.Context) (int, error)
	CreateGroup(ctx context.Context, group *models.Group) (*models.Group, error)
	UpdateGroup(ctx context.Context, group *models.Group) (*models.Group, error)
	DeleteGroupById(ctx context.Context, id uuid.UUID) error
	GetGroupById(ctx context.Context, id uuid.UUID) (*models.Group, error)
	CountGroups(ctx context.Context) (int, error)
	CreateMembership(ctx context.Context, membership *models.Membership) (*models.Membership, error)
	UpdateMembership(ctx context.Context, membership *models.Membership) (*models.Membership, error)
	DeleteMembershipById(ctx context.Context, id uuid.UUID) error
	GetMembershipById(ctx context.Context, id uuid.UUID) (*models.Membership, error)
	CountMemberships(ctx context.Context) (int, error)
	GetUserMembershipById(ctx context.Context, id uuid.UUID) (*models.UserMembership, error)
	GetGroupMembershipById(ctx context.Context, id uuid.UUID) (*models.GroupMembership, error)
	BlacklistToken(ctx context.Context, blacklist *models.Blacklist) (*models.Blacklist, error)
	CheckBlacklist(ctx context.Context, accessToken string) (*models.Blacklist, error)
	UpdateUserPassword(ctx context.Context, user *models.User) (*models.User, error)
}
