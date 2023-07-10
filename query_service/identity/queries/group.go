package queries

import (
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/gofrs/uuid"
)

type GroupQueries struct {
	GetGroupById GetGroupByIdHandler
	SearchGroup  SearchGroupHandler
}

func NewGroupQueries(getById GetGroupByIdHandler, search SearchGroupHandler) *GroupQueries {
	return &GroupQueries{
		GetGroupById: getById,
		SearchGroup:  search,
	}
}

type GetGroupByIdQuery struct {
	ID uuid.UUID `json:"id" bson:"_id,omitempty"`
}

func NewGetGroupByIdQuery(id uuid.UUID) *GetGroupByIdQuery {
	return &GetGroupByIdQuery{ID: id}
}

type SearchGroupQuery struct {
	Text       string                `json:"text"`
	Pagination *utilities.Pagination `json:"pagination"`
}

func NewSearchGroupQuery(text string, pagination *utilities.Pagination) *SearchGroupQuery {
	return &SearchGroupQuery{
		Text:       text,
		Pagination: pagination,
	}
}
