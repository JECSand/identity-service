package queries

import (
	"github.com/gofrs/uuid"
)

// UserQueries ...
type UserQueries struct {
	GetUserById GetUserByIdHandler
	CountUsers  CountUsersHandler
}

// NewUserQueries ...
func NewUserQueries(getById GetUserByIdHandler, countUsers CountUsersHandler) *UserQueries {
	return &UserQueries{
		GetUserById: getById,
		CountUsers:  countUsers,
	}
}

// GetUserByIdQuery ...
type GetUserByIdQuery struct {
	ID uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
}

// NewGetUserByIdQuery ...
func NewGetUserByIdQuery(id uuid.UUID) *GetUserByIdQuery {
	return &GetUserByIdQuery{ID: id}
}
