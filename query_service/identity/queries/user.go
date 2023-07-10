package queries

import (
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/gofrs/uuid"
)

type UserQueries struct {
	GetUserById GetUserByIdHandler
	SearchUser  SearchUserHandler
}

func NewUserQueries(getById GetUserByIdHandler, search SearchUserHandler) *UserQueries {
	return &UserQueries{
		GetUserById: getById,
		SearchUser:  search,
	}
}

type GetUserByIdQuery struct {
	ID uuid.UUID `json:"id" bson:"_id,omitempty"`
}

func NewGetUserByIdQuery(id uuid.UUID) *GetUserByIdQuery {
	return &GetUserByIdQuery{ID: id}
}

type SearchUserQuery struct {
	Text       string                `json:"text"`
	Pagination *utilities.Pagination `json:"pagination"`
}

func NewSearchUserQuery(text string, pagination *utilities.Pagination) *SearchUserQuery {
	return &SearchUserQuery{
		Text:       text,
		Pagination: pagination,
	}
}
