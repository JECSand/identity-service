package cache

import (
	"context"
	"github.com/JECSand/identity-service/query_service/identity/entities"
)

type Cache interface {
	PutUserMembership(ctx context.Context, key string, userMembership *entities.UserMembership)
	GetUserMembership(ctx context.Context, key string) (*entities.UserMembership, error)
	DeleteUserMembership(ctx context.Context, key string)
	DeleteAllUserMemberships(ctx context.Context)
	PutGroupMembership(ctx context.Context, key string, groupMembership *entities.GroupMembership)
	GetGroupMembership(ctx context.Context, key string) (*entities.GroupMembership, error)
	DeleteGroupMembership(ctx context.Context, key string)
	DeleteAllGroupMemberships(ctx context.Context)
	PutMembership(ctx context.Context, key string, membership *entities.Membership)
	GetMembership(ctx context.Context, key string) (*entities.Membership, error)
	DeleteMembership(ctx context.Context, key string)
	DeleteAllMemberships(ctx context.Context)
	PutGroup(ctx context.Context, key string, group *entities.Group)
	GetGroup(ctx context.Context, key string) (*entities.Group, error)
	DeleteGroup(ctx context.Context, key string)
	DeleteAllGroups(ctx context.Context)
	PutUser(ctx context.Context, key string, user *entities.User)
	GetUser(ctx context.Context, key string) (*entities.User, error)
	DeleteUser(ctx context.Context, key string)
	DeleteAllUsers(ctx context.Context)
	PutToken(ctx context.Context, key string, blacklist *entities.Blacklist)
	GetToken(ctx context.Context, key string) (*entities.Blacklist, error)
	DeleteToken(ctx context.Context, key string)
	DeleteAllTokens(ctx context.Context)
}
