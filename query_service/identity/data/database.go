package data

import (
	"context"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// database structures the available mongo services
type database struct {
	users            *userRepository
	groups           *groupRepository
	memberships      *membershipRepository
	groupMemberships *groupMembershipRepository
	userMemberships  *userMembershipRepository
	blacklist        *blacklistRepository
}

// NewDatabase Initializes a new Database setup to MongoDB
func NewDatabase(log logging.Logger, cfg *config.Config, db *mongo.Client) *database {
	userRepo := NewUserRepository(log, cfg, db)
	groupRepo := NewGroupRepository(log, cfg, db)
	membershipRepo := NewMembershipRepository(log, cfg, db)
	groupMembershipRepo := NewGroupMembershipRepository(log, cfg, db)
	userMembershipRepo := NewUserMembershipRepository(log, cfg, db)
	blRepo := NewBlacklistRepository(log, cfg, db)
	return &database{
		userRepo,
		groupRepo,
		membershipRepo,
		groupMembershipRepo,
		userMembershipRepo,
		blRepo,
	}
}

func (d *database) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	return d.users.Create(ctx, user)
}

func (d *database) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	return d.users.Update(ctx, user)
}

func (d *database) GetUserById(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return d.users.GetById(ctx, id)
}

func (d *database) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return d.users.GetByEmail(ctx, email)
}

func (d *database) AuthenticateUser(ctx context.Context, email string, password string) (*entities.User, error) {
	return d.users.Authenticate(ctx, email, password)
}

func (d *database) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return d.users.Delete(ctx, id)
}

func (d *database) SearchUsers(ctx context.Context, search string, pagination *utilities.Pagination) (*entities.UsersList, error) {
	return d.users.Search(ctx, search, pagination)
}

func (d *database) CreateGroup(ctx context.Context, model *entities.Group) (*entities.Group, error) {
	return d.groups.Create(ctx, model)
}

func (d *database) UpdateGroup(ctx context.Context, model *entities.Group) (*entities.Group, error) {
	return d.groups.Update(ctx, model)
}

func (d *database) GetGroupById(ctx context.Context, id uuid.UUID) (*entities.Group, error) {
	return d.groups.GetById(ctx, id)
}

func (d *database) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	return d.groups.Delete(ctx, id)
}

func (d *database) SearchGroups(ctx context.Context, search string, pagination *utilities.Pagination) (*entities.GroupsList, error) {
	return d.groups.Search(ctx, search, pagination)
}

func (d *database) CreateMembership(ctx context.Context, model *entities.Membership) (*entities.Membership, error) {
	return d.memberships.Create(ctx, model)
}

func (d *database) UpdateMembership(ctx context.Context, model *entities.Membership) (*entities.Membership, error) {
	return d.memberships.Update(ctx, model)
}

func (d *database) UpdateMemberships(ctx context.Context, filter *entities.Membership, update *entities.Membership) error {
	return d.memberships.UpdateMany(ctx, filter, update)
}

func (d *database) GetMembershipById(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	return d.memberships.GetById(ctx, id)
}

func (d *database) DeleteMembership(ctx context.Context, id uuid.UUID) error {
	return d.memberships.Delete(ctx, id)
}

func (d *database) DeleteMemberships(ctx context.Context, filter *entities.Membership) error {
	return d.memberships.DeleteMany(ctx, filter)
}

func (d *database) CreateGroupMembership(ctx context.Context, model *entities.GroupMembership) (*entities.GroupMembership, error) {
	return d.groupMemberships.Create(ctx, model)
}

func (d *database) UpdateGroupMembership(ctx context.Context, model *entities.GroupMembership) (*entities.GroupMembership, error) {
	return d.groupMemberships.Update(ctx, model)
}

func (d *database) UpdateGroupMemberships(ctx context.Context, filter *entities.GroupMembership, update *entities.GroupMembership) error {
	return d.groupMemberships.UpdateMany(ctx, filter, update)
}

func (d *database) GetGroupMembershipById(ctx context.Context, id uuid.UUID, idType enums.ReadTableIdType) (*entities.GroupMembership, error) {
	return d.groupMemberships.GetById(ctx, id, idType)
}

func (d *database) GetGroupMembershipByUserId(ctx context.Context, userId uuid.UUID, pagination *utilities.Pagination) (*entities.GroupMembershipsList, error) {
	return d.groupMemberships.GetByUserId(ctx, userId, pagination)
}

func (d *database) GetGroupMembershipByGroupId(ctx context.Context, groupId uuid.UUID, pagination *utilities.Pagination) (*entities.GroupMembershipsList, error) {
	return d.groupMemberships.GetByGroupId(ctx, groupId, pagination)
}

func (d *database) DeleteGroupMembership(ctx context.Context, id uuid.UUID) error {
	return d.groupMemberships.Delete(ctx, id)
}

func (d *database) DeleteGroupMemberships(ctx context.Context, filter *entities.GroupMembership) error {
	return d.groupMemberships.DeleteMany(ctx, filter)
}

func (d *database) DeleteGroupMembershipByMembershipId(ctx context.Context, id uuid.UUID) error {
	return d.userMemberships.DeleteByMembershipId(ctx, id)
}

func (d *database) CreateUserMembership(ctx context.Context, model *entities.UserMembership) (*entities.UserMembership, error) {
	return d.userMemberships.Create(ctx, model)
}

func (d *database) UpdateUserMembership(ctx context.Context, model *entities.UserMembership) (*entities.UserMembership, error) {
	return d.userMemberships.Update(ctx, model)
}

func (d *database) UpdateUserMemberships(ctx context.Context, filter *entities.UserMembership, update *entities.UserMembership) error {
	return d.userMemberships.UpdateMany(ctx, filter, update)
}

func (d *database) GetUserMembershipById(ctx context.Context, id uuid.UUID, idType enums.ReadTableIdType) (*entities.UserMembership, error) {
	return d.userMemberships.GetById(ctx, id, idType)
}

func (d *database) GetUserMembershipByUserId(ctx context.Context, userId uuid.UUID, pagination *utilities.Pagination) (*entities.UserMembershipsList, error) {
	return d.userMemberships.GetByUserId(ctx, userId, pagination)
}

func (d *database) GetUserMembershipByGroupId(ctx context.Context, groupId uuid.UUID, pagination *utilities.Pagination) (*entities.UserMembershipsList, error) {
	return d.userMemberships.GetByGroupId(ctx, groupId, pagination)
}

func (d *database) DeleteUserMembership(ctx context.Context, id uuid.UUID) error {
	return d.userMemberships.Delete(ctx, id)
}

func (d *database) DeleteUserMemberships(ctx context.Context, filter *entities.UserMembership) error {
	return d.userMemberships.DeleteMany(ctx, filter)
}

func (d *database) DeleteUserMembershipByMembershipId(ctx context.Context, id uuid.UUID) error {
	return d.userMemberships.DeleteByMembershipId(ctx, id)
}

func (d *database) BlacklistToken(ctx context.Context, bList *entities.Blacklist) (*entities.Blacklist, error) {
	return d.blacklist.BlacklistToken(ctx, bList)
}

func (d *database) CheckTokenBlacklist(ctx context.Context, accessToken string) (*entities.Blacklist, error) {
	return d.blacklist.CheckBlacklist(ctx, accessToken)
}

type Database interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	AuthenticateUser(ctx context.Context, email string, password string) (*entities.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SearchUsers(ctx context.Context, search string, pagination *utilities.Pagination) (*entities.UsersList, error)
	CreateGroup(ctx context.Context, model *entities.Group) (*entities.Group, error)
	UpdateGroup(ctx context.Context, model *entities.Group) (*entities.Group, error)
	GetGroupById(ctx context.Context, id uuid.UUID) (*entities.Group, error)
	DeleteGroup(ctx context.Context, id uuid.UUID) error
	SearchGroups(ctx context.Context, search string, pagination *utilities.Pagination) (*entities.GroupsList, error)
	CreateMembership(ctx context.Context, model *entities.Membership) (*entities.Membership, error)
	UpdateMembership(ctx context.Context, model *entities.Membership) (*entities.Membership, error)
	UpdateMemberships(ctx context.Context, filter *entities.Membership, update *entities.Membership) error
	GetMembershipById(ctx context.Context, id uuid.UUID) (*entities.Membership, error)
	DeleteMembership(ctx context.Context, id uuid.UUID) error
	DeleteMemberships(ctx context.Context, filter *entities.Membership) error
	CreateGroupMembership(ctx context.Context, model *entities.GroupMembership) (*entities.GroupMembership, error)
	UpdateGroupMembership(ctx context.Context, model *entities.GroupMembership) (*entities.GroupMembership, error)
	UpdateGroupMemberships(ctx context.Context, filter *entities.GroupMembership, update *entities.GroupMembership) error
	GetGroupMembershipById(ctx context.Context, id uuid.UUID, idType enums.ReadTableIdType) (*entities.GroupMembership, error)
	GetGroupMembershipByUserId(ctx context.Context, userId uuid.UUID, pagination *utilities.Pagination) (*entities.GroupMembershipsList, error)
	GetGroupMembershipByGroupId(ctx context.Context, groupId uuid.UUID, pagination *utilities.Pagination) (*entities.GroupMembershipsList, error)
	DeleteGroupMembership(ctx context.Context, id uuid.UUID) error
	DeleteGroupMemberships(ctx context.Context, filter *entities.GroupMembership) error
	DeleteGroupMembershipByMembershipId(ctx context.Context, id uuid.UUID) error
	CreateUserMembership(ctx context.Context, model *entities.UserMembership) (*entities.UserMembership, error)
	UpdateUserMembership(ctx context.Context, model *entities.UserMembership) (*entities.UserMembership, error)
	UpdateUserMemberships(ctx context.Context, filter *entities.UserMembership, update *entities.UserMembership) error
	GetUserMembershipById(ctx context.Context, id uuid.UUID, idType enums.ReadTableIdType) (*entities.UserMembership, error)
	GetUserMembershipByUserId(ctx context.Context, userId uuid.UUID, pagination *utilities.Pagination) (*entities.UserMembershipsList, error)
	GetUserMembershipByGroupId(ctx context.Context, groupId uuid.UUID, pagination *utilities.Pagination) (*entities.UserMembershipsList, error)
	DeleteUserMembership(ctx context.Context, id uuid.UUID) error
	DeleteUserMemberships(ctx context.Context, filter *entities.UserMembership) error
	DeleteUserMembershipByMembershipId(ctx context.Context, id uuid.UUID) error
	BlacklistToken(ctx context.Context, bList *entities.Blacklist) (*entities.Blacklist, error)
	CheckTokenBlacklist(ctx context.Context, accessToken string) (*entities.Blacklist, error)
}
