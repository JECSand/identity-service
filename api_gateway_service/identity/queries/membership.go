package queries

import (
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/gofrs/uuid"
)

type MembershipQueries struct {
	GetMembershipById          GetMembershipByIdHandler
	GetUserMembershipByGroupId GetUserMembershipByGroupIdHandler
	GetGroupMembershipByUserId GetGroupMembershipByUserIdHandler
}

func NewMembershipQueries(getById GetMembershipByIdHandler, getUserMembership GetUserMembershipByGroupIdHandler, getGroupMembership GetGroupMembershipByUserIdHandler) *MembershipQueries {
	return &MembershipQueries{
		GetMembershipById:          getById,
		GetUserMembershipByGroupId: getUserMembership,
		GetGroupMembershipByUserId: getGroupMembership,
	}
}

type GetMembershipByIdQuery struct {
	ID uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
}

func NewGetMembershipByIdQuery(id uuid.UUID) *GetMembershipByIdQuery {
	return &GetMembershipByIdQuery{ID: id}
}

type GetUserMembershipByGroupIdQuery struct {
	GroupID    uuid.UUID             `json:"groupID" validate:"required,gte=0,lte=255"`
	Pagination *utilities.Pagination `json:"pagination"`
}

func NewGetUserMembershipByGroupIdQuery(groupId uuid.UUID, pagination *utilities.Pagination) *GetUserMembershipByGroupIdQuery {
	return &GetUserMembershipByGroupIdQuery{
		GroupID:    groupId,
		Pagination: pagination,
	}
}

type GetGroupMembershipByUserIdQuery struct {
	UserID     uuid.UUID             `json:"userID" validate:"required,gte=0,lte=255"`
	Pagination *utilities.Pagination `json:"pagination"`
}

func NewGetGroupMembershipByUserIdQuery(userId uuid.UUID, pagination *utilities.Pagination) *GetGroupMembershipByUserIdQuery {
	return &GetGroupMembershipByUserIdQuery{
		UserID:     userId,
		Pagination: pagination,
	}
}
