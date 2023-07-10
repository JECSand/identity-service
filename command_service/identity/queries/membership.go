package queries

import (
	"github.com/gofrs/uuid"
)

// MembershipQueries ...
type MembershipQueries struct {
	GetMembershipById      GetMembershipByIdHandler
	GetUserMembershipById  GetUserMembershipByIdHandler
	GetGroupMembershipById GetGroupMembershipByIdHandler
	CountMemberships       CountMembershipsHandler
}

// NewMembershipQueries ...
func NewMembershipQueries(
	getById GetMembershipByIdHandler,
	getUserMembershipById GetUserMembershipByIdHandler,
	getGroupMembershipById GetGroupMembershipByIdHandler,
	countMemberships CountMembershipsHandler,
) *MembershipQueries {
	return &MembershipQueries{
		GetMembershipById:      getById,
		GetUserMembershipById:  getUserMembershipById,
		GetGroupMembershipById: getGroupMembershipById,
		CountMemberships:       countMemberships,
	}
}

// GetMembershipByIdQuery ...
type GetMembershipByIdQuery struct {
	ID uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
}

// NewGetMembershipByIdQuery ...
func NewGetMembershipByIdQuery(id uuid.UUID) *GetMembershipByIdQuery {
	return &GetMembershipByIdQuery{ID: id}
}

// GetUserMembershipByIdQuery ...
type GetUserMembershipByIdQuery struct {
	ID uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
}

// NewGetUserMembershipByIdQuery ...
func NewGetUserMembershipByIdQuery(id uuid.UUID) *GetUserMembershipByIdQuery {
	return &GetUserMembershipByIdQuery{ID: id}
}

// GetGroupMembershipByIdQuery ...
type GetGroupMembershipByIdQuery struct {
	ID uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
}

// NewGetGroupMembershipByIdQuery ...
func NewGetGroupMembershipByIdQuery(id uuid.UUID) *GetGroupMembershipByIdQuery {
	return &GetGroupMembershipByIdQuery{ID: id}
}
