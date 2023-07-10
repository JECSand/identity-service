package queries

import (
	"github.com/gofrs/uuid"
)

// GroupQueries ...
type GroupQueries struct {
	GetGroupById GetGroupByIdHandler
	CountGroups  CountGroupsHandler
}

// NewGroupQueries ...
func NewGroupQueries(getById GetGroupByIdHandler, countGroups CountGroupsHandler) *GroupQueries {
	return &GroupQueries{
		GetGroupById: getById,
		CountGroups:  countGroups,
	}
}

// GetGroupByIdQuery ...
type GetGroupByIdQuery struct {
	ID uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
}

// NewGetGroupByIdQuery ...
func NewGetGroupByIdQuery(id uuid.UUID) *GetGroupByIdQuery {
	return &GetGroupByIdQuery{ID: id}
}
