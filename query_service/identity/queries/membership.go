package queries

import (
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/gofrs/uuid"
)

type MembershipQueries struct {
	GetMembershipById  GetMembershipByIdHandler
	GetGroupMembership GetGroupMembershipHandler
	GetUserMembership  GetUserMembershipHandler
}

func NewMembershipQueries(getById GetMembershipByIdHandler, userGroups GetGroupMembershipHandler, groupUsers GetUserMembershipHandler) *MembershipQueries {
	return &MembershipQueries{
		GetMembershipById:  getById,
		GetGroupMembership: userGroups,
		GetUserMembership:  groupUsers,
	}
}

type GetMembershipByIdQuery struct {
	ID uuid.UUID `json:"id" bson:"_id,omitempty"`
}

func NewGetMembershipByIdQuery(id uuid.UUID) *GetMembershipByIdQuery {
	return &GetMembershipByIdQuery{ID: id}
}

type GetGroupMembershipQuery struct {
	UserID     uuid.UUID             `json:"userID"`
	Pagination *utilities.Pagination `json:"pagination"`
}

func NewGetGroupMembershipQuery(userID uuid.UUID, pagination *utilities.Pagination) *GetGroupMembershipQuery {
	return &GetGroupMembershipQuery{
		UserID:     userID,
		Pagination: pagination,
	}
}

type GetUserMembershipQuery struct {
	GroupID    uuid.UUID             `json:"groupID"`
	Pagination *utilities.Pagination `json:"pagination"`
}

func NewGetUserMembershipQuery(groupID uuid.UUID, pagination *utilities.Pagination) *GetUserMembershipQuery {
	return &GetUserMembershipQuery{
		GroupID:    groupID,
		Pagination: pagination,
	}
}
