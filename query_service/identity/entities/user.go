package entities

import (
	"errors"
	"github.com/JECSand/identity-service/pkg/utilities"
	authQueryService "github.com/JECSand/identity-service/query_service/protos/auth_query"
	queryService "github.com/JECSand/identity-service/query_service/protos/user_query"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty" validate:"required,min=3,max=250"`
	Username  string    `json:"username,omitempty" bson:"username,omitempty" validate:"required,min=3,max=500"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	Root      bool      `json:"root,omitempty" bson:"root,omitempty"`
	Active    bool      `json:"active,omitempty" bson:"active,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// GetID returns the unique identifier of the User
func (u *User) GetID() string {
	return u.ID
}

// Authenticate compares an input password with the hashed password stored in the User model
func (u *User) Authenticate(checkPassword string) error {
	if len(u.Password) != 0 {
		password := []byte(u.Password)
		cPassword := []byte(checkPassword)
		return bcrypt.CompareHashAndPassword(password, cPassword)
	}
	return errors.New("user password is missing")
}

// UsersList response with pagination
type UsersList struct {
	TotalCount int64   `json:"totalCount" bson:"totalCount"`
	TotalPages int64   `json:"totalPages" bson:"totalPages"`
	Page       int64   `json:"page" bson:"page"`
	Size       int64   `json:"size" bson:"size"`
	HasMore    bool    `json:"hasMore" bson:"hasMore"`
	Users      []*User `json:"users" bson:"users"`
}

func NewUserListWithPagination(users []*User, count int64, pagination *utilities.Pagination) *UsersList {
	return &UsersList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Users:      users,
	}
}

func AuthUserToGrpcMessage(user *User) *authQueryService.User {
	return &authQueryService.User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		Root:      user.Root,
		Active:    user.Active,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func UserToGrpcMessage(user *User) *queryService.User {
	return &queryService.User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		Root:      user.Root,
		Active:    user.Active,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func UserListToGrpc(users *UsersList) *queryService.SearchRes {
	list := make([]*queryService.User, 0, len(users.Users))
	for _, user := range users.Users {
		list = append(list, UserToGrpcMessage(user))
	}
	return &queryService.SearchRes{
		TotalCount: users.TotalCount,
		TotalPages: users.TotalPages,
		Page:       users.Page,
		Size:       users.Size,
		HasMore:    users.HasMore,
		Users:      list,
	}
}
