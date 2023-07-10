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
	ID uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
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
