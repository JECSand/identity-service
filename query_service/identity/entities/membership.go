package entities

import (
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/utilities"
	queryService "github.com/JECSand/identity-service/query_service/protos/membership_query"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Membership struct {
	ID        string                 `json:"id" bson:"_id,omitempty"`
	UserID    string                 `json:"userID,omitempty" bson:"user_id,omitempty" validate:"required,min=3,max=250"`
	GroupID   string                 `json:"groupID,omitempty" bson:"group_id,omitempty" validate:"required,min=3,max=500"`
	Status    enums.MembershipStatus `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
	Role      enums.Role             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt time.Time              `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time              `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// GetID returns the unique identifier of the Membership
func (u *Membership) GetID() string {
	return u.ID
}

func MembershipToGrpcMessage(membership *Membership) *queryService.Membership {
	return &queryService.Membership{
		ID:        membership.ID,
		UserID:    membership.UserID,
		GroupID:   membership.GroupID,
		Status:    int64(membership.Status.EnumIndex()),
		Role:      int64(membership.Role.EnumIndex()),
		CreatedAt: timestamppb.New(membership.CreatedAt),
		UpdatedAt: timestamppb.New(membership.UpdatedAt),
	}
}

type UserMembership struct {
	ID           string                 `json:"id" bson:"_id,omitempty"`
	GroupID      string                 `json:"groupID,omitempty" bson:"group_id,omitempty" validate:"required,min=3,max=500"`
	UserID       string                 `json:"userID,omitempty" bson:"user_id,omitempty" validate:"required,min=3,max=250"`
	MembershipID string                 `json:"membershipID,omitempty" bson:"membership_id,omitempty" validate:"required,min=3,max=500"`
	Email        string                 `json:"email,omitempty" bson:"email,omitempty" validate:"required,min=3,max=500"`
	Username     string                 `json:"username,omitempty" bson:"username,omitempty" validate:"required,min=3,max=500"`
	Status       enums.MembershipStatus `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
	Role         enums.Role             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt    time.Time              `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// GetID returns the unique identifier of the UserMembership
func (u *UserMembership) GetID() string {
	return u.ID
}

// UserMembershipsList response with pagination
type UserMembershipsList struct {
	TotalCount      int64             `json:"totalCount" bson:"total_count"`
	TotalPages      int64             `json:"totalPages" bson:"total_pages"`
	Page            int64             `json:"page" bson:"page"`
	Size            int64             `json:"size" bson:"size"`
	HasMore         bool              `json:"hasMore" bson:"has_more"`
	UserMemberships []*UserMembership `json:"userMemberships" bson:"user_memberships"`
}

func UserMembershipToGrpcMessage(userMembership *UserMembership) *queryService.UserMembership {
	return &queryService.UserMembership{
		ID:           userMembership.ID,
		GroupID:      userMembership.GroupID,
		UserID:       userMembership.UserID,
		MembershipID: userMembership.MembershipID,
		Email:        userMembership.Email,
		Username:     userMembership.Username,
		Status:       int64(userMembership.Status.EnumIndex()),
		Role:         int64(userMembership.Role.EnumIndex()),
		CreatedAt:    timestamppb.New(userMembership.CreatedAt),
		UpdatedAt:    timestamppb.New(userMembership.UpdatedAt),
	}
}

func NewUserMembershipListWithPagination(memberships []*UserMembership, count int64, pagination *utilities.Pagination) *UserMembershipsList {
	return &UserMembershipsList{
		TotalCount:      count,
		TotalPages:      int64(pagination.GetTotalPages(int(count))),
		Page:            int64(pagination.GetPage()),
		Size:            int64(pagination.GetSize()),
		HasMore:         pagination.GetHasMore(int(count)),
		UserMemberships: memberships,
	}
}

func UserMembershipListToGrpc(userMemberships *UserMembershipsList) *queryService.GetUserMembershipRes {
	list := make([]*queryService.UserMembership, 0, len(userMemberships.UserMemberships))
	for _, userMember := range userMemberships.UserMemberships {
		list = append(list, UserMembershipToGrpcMessage(userMember))
	}
	return &queryService.GetUserMembershipRes{
		TotalCount:      userMemberships.TotalCount,
		TotalPages:      userMemberships.TotalPages,
		Page:            userMemberships.Page,
		Size:            userMemberships.Size,
		HasMore:         userMemberships.HasMore,
		UserMemberships: list,
	}
}

type GroupMembership struct {
	ID           string                 `json:"id" bson:"_id,omitempty"`
	UserID       string                 `json:"userID,omitempty" bson:"user_id,omitempty" validate:"required,min=3,max=250"`
	GroupID      string                 `json:"groupID,omitempty" bson:"group_id,omitempty" validate:"required,min=3,max=500"`
	MembershipID string                 `json:"membershipID,omitempty" bson:"membership_id,omitempty" validate:"required,min=3,max=500"`
	Name         string                 `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=500"`
	Description  string                 `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Status       enums.MembershipStatus `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
	Role         enums.Role             `json:"role,omitempty" bson:"role,omitempty"`
	Creator      bool                   `json:"creator,omitempty" bson:"creator,omitempty"`
	CreatedAt    time.Time              `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// GetID returns the unique identifier of the GroupMembership
func (u *GroupMembership) GetID() string {
	return u.ID
}

// GroupMembershipsList response with pagination
type GroupMembershipsList struct {
	TotalCount       int64              `json:"totalCount" bson:"total_count"`
	TotalPages       int64              `json:"totalPages" bson:"total_pages"`
	Page             int64              `json:"page" bson:"page"`
	Size             int64              `json:"size" bson:"size"`
	HasMore          bool               `json:"hasMore" bson:"has_more"`
	GroupMemberships []*GroupMembership `json:"groupMemberships" bson:"group_memberships"`
}

func GroupMembershipToGrpcMessage(groupMembership *GroupMembership) *queryService.GroupMembership {
	return &queryService.GroupMembership{
		ID:           groupMembership.ID,
		UserID:       groupMembership.UserID,
		GroupID:      groupMembership.GroupID,
		MembershipID: groupMembership.MembershipID,
		Name:         groupMembership.Name,
		Description:  groupMembership.Description,
		Status:       int64(groupMembership.Status.EnumIndex()),
		Role:         int64(groupMembership.Role.EnumIndex()),
		CreatedAt:    timestamppb.New(groupMembership.CreatedAt),
		UpdatedAt:    timestamppb.New(groupMembership.UpdatedAt),
	}
}

func NewGroupMembershipListWithPagination(memberships []*GroupMembership, count int64, pagination *utilities.Pagination) *GroupMembershipsList {
	return &GroupMembershipsList{
		TotalCount:       count,
		TotalPages:       int64(pagination.GetTotalPages(int(count))),
		Page:             int64(pagination.GetPage()),
		Size:             int64(pagination.GetSize()),
		HasMore:          pagination.GetHasMore(int(count)),
		GroupMemberships: memberships,
	}
}

func GroupMembershipListToGrpc(groupMemberships *GroupMembershipsList) *queryService.GetGroupMembershipRes {
	list := make([]*queryService.GroupMembership, 0, len(groupMemberships.GroupMemberships))
	for _, group := range groupMemberships.GroupMemberships {
		list = append(list, GroupMembershipToGrpcMessage(group))
	}
	return &queryService.GetGroupMembershipRes{
		TotalCount:       groupMemberships.TotalCount,
		TotalPages:       groupMemberships.TotalPages,
		Page:             groupMemberships.Page,
		Size:             groupMemberships.Size,
		HasMore:          groupMemberships.HasMore,
		GroupMemberships: list,
	}
}
