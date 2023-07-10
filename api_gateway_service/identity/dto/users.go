package dto

import (
	queryService "github.com/JECSand/identity-service/query_service/protos/user_query"
	"github.com/gofrs/uuid"
	"time"
)

type CreateUserDTO struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email" validate:"required,gte=0,lte=255"`
	Username string    `json:"username" validate:"required,gte=0,lte=255"`
	Password string    `json:"password" validate:"required,gte=0,lte=5000"`
	Active   bool      `json:"active"`
}

type CreateUserResponseDTO struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type UpdateUserDTO struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Email    string    `json:"email" validate:"required,gte=0,lte=255"`
	Username string    `json:"username" validate:"required,gte=0,lte=255"`
	Password string    `json:"password" validate:"required,gte=0,lte=5000"`
	Active   bool      `json:"active"`
}

// UserResponse ...
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email,omitempty"`
	Username  string    `json:"username,omitempty"`
	Root      bool      `json:"root,omitempty"`
	Active    bool      `json:"active,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func UserResponseFromGrpc(user *queryService.User) *UserResponse {
	return &UserResponse{
		ID:        user.GetID(),
		Email:     user.GetEmail(),
		Username:  user.GetUsername(),
		Root:      user.GetRoot(),
		Active:    user.GetActive(),
		CreatedAt: user.GetCreatedAt().AsTime(),
		UpdatedAt: user.GetUpdatedAt().AsTime(),
	}
}

// UsersListResponse ...
type UsersListResponse struct {
	TotalCount int64           `json:"totalCount" bson:"total_count"`
	TotalPages int64           `json:"totalPages" bson:"total_pages"`
	Page       int64           `json:"page" bson:"page"`
	Size       int64           `json:"size" bson:"size"`
	HasMore    bool            `json:"hasMore" bson:"has_more"`
	Users      []*UserResponse `json:"users" bson:"users"`
}

func UsersListResponseFromGrpc(listResponse *queryService.SearchRes) *UsersListResponse {
	list := make([]*UserResponse, 0, len(listResponse.GetUsers()))
	for _, user := range listResponse.GetUsers() {
		list = append(list, UserResponseFromGrpc(user))
	}
	return &UsersListResponse{
		TotalCount: listResponse.GetTotalCount(),
		TotalPages: listResponse.GetTotalPages(),
		Page:       listResponse.GetPage(),
		Size:       listResponse.GetSize(),
		HasMore:    listResponse.GetHasMore(),
		Users:      list,
	}
}
