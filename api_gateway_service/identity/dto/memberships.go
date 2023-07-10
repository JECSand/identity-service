package dto

import (
	"github.com/JECSand/identity-service/pkg/enums"
	membershipQueryService "github.com/JECSand/identity-service/query_service/protos/membership_query"
	"github.com/gofrs/uuid"
	"time"
)

type CreateMembershipDTO struct {
	ID      uuid.UUID              `json:"id"`
	UserID  uuid.UUID              `json:"userID" validate:"required"`
	GroupID uuid.UUID              `json:"groupID" validate:"required"`
	Status  enums.MembershipStatus `json:"status" validate:"required"`
	Role    enums.Role             `json:"role" validate:"required"`
}

type CreateMembershipResponseDTO struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type UpdateMembershipDTO struct {
	ID     uuid.UUID              `json:"id" validate:"required"`
	Status enums.MembershipStatus `json:"status" validate:"required"`
	Role   enums.Role             `json:"role" validate:"required"`
}

// MembershipResponse ...
type MembershipResponse struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"userID,omitempty"`
	GroupID   string                 `json:"groupID,omitempty"`
	Status    enums.MembershipStatus `json:"status,omitempty"`
	Role      enums.Role             `json:"role,omitempty"`
	CreatedAt time.Time              `json:"createdAt,omitempty"`
	UpdatedAt time.Time              `json:"updatedAt,omitempty"`
}

func MembershipResponseFromGrpc(membership *membershipQueryService.Membership) *MembershipResponse {
	return &MembershipResponse{
		ID:        membership.GetID(),
		UserID:    membership.GetUserID(),
		GroupID:   membership.GetGroupID(),
		Status:    enums.MembershipStatus(membership.GetStatus()),
		Role:      enums.Role(membership.GetRole()),
		CreatedAt: membership.GetCreatedAt().AsTime(),
		UpdatedAt: membership.GetUpdatedAt().AsTime(),
	}
}

// UserMembershipResponse ...
type UserMembershipResponse struct {
	ID           string                 `json:"id"`
	GroupID      string                 `json:"groupID,omitempty"`
	UserID       string                 `json:"userID,omitempty"`
	MembershipID string                 `json:"membershipID,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Username     string                 `json:"username,omitempty"`
	Status       enums.MembershipStatus `json:"status,omitempty"`
	Role         enums.Role             `json:"role,omitempty"`
	CreatedAt    time.Time              `json:"createdAt,omitempty"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty"`
}

func UserMembershipResponseFromGrpc(userMembership *membershipQueryService.UserMembership) *UserMembershipResponse {
	return &UserMembershipResponse{
		ID:           userMembership.GetID(),
		GroupID:      userMembership.GetGroupID(),
		UserID:       userMembership.GetUserID(),
		MembershipID: userMembership.GetMembershipID(),
		Email:        userMembership.GetEmail(),
		Username:     userMembership.GetUsername(),
		Status:       enums.MembershipStatus(userMembership.GetStatus()),
		Role:         enums.Role(userMembership.GetRole()),
		CreatedAt:    userMembership.GetCreatedAt().AsTime(),
		UpdatedAt:    userMembership.GetUpdatedAt().AsTime(),
	}
}

// UserMembershipsListResponse ...
type UserMembershipsListResponse struct {
	TotalCount      int64                     `json:"totalCount" bson:"total_count"`
	TotalPages      int64                     `json:"totalPages" bson:"total_pages"`
	Page            int64                     `json:"page" bson:"page"`
	Size            int64                     `json:"size" bson:"size"`
	HasMore         bool                      `json:"hasMore" bson:"has_more"`
	UserMemberships []*UserMembershipResponse `json:"userMemberships" bson:"user_memberships"`
}

func UserMembershipListResponseFromGrpc(listResponse *membershipQueryService.GetUserMembershipRes) *UserMembershipsListResponse {
	list := make([]*UserMembershipResponse, 0, len(listResponse.GetUserMemberships()))
	for _, um := range listResponse.GetUserMemberships() {
		list = append(list, UserMembershipResponseFromGrpc(um))
	}
	return &UserMembershipsListResponse{
		TotalCount:      listResponse.GetTotalCount(),
		TotalPages:      listResponse.GetTotalPages(),
		Page:            listResponse.GetPage(),
		Size:            listResponse.GetSize(),
		HasMore:         listResponse.GetHasMore(),
		UserMemberships: list,
	}
}

// GroupMembershipResponse ...
type GroupMembershipResponse struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"userID,omitempty"`
	GroupID      string                 `json:"groupID,omitempty"`
	MembershipID string                 `json:"membershipID,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Status       enums.MembershipStatus `json:"status,omitempty"`
	Role         enums.Role             `json:"role,omitempty"`
	Creator      bool                   `json:"creator,omitempty"`
	CreatedAt    time.Time              `json:"createdAt,omitempty"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty"`
}

func GroupMembershipResponseFromGrpc(groupMembership *membershipQueryService.GroupMembership) *GroupMembershipResponse {
	return &GroupMembershipResponse{
		ID:           groupMembership.GetID(),
		UserID:       groupMembership.GetUserID(),
		GroupID:      groupMembership.GetGroupID(),
		MembershipID: groupMembership.GetMembershipID(),
		Name:         groupMembership.GetName(),
		Description:  groupMembership.GetDescription(),
		Status:       enums.MembershipStatus(groupMembership.GetStatus()),
		Role:         enums.Role(groupMembership.GetRole()),
		Creator:      groupMembership.GetCreator(),
		CreatedAt:    groupMembership.GetCreatedAt().AsTime(),
		UpdatedAt:    groupMembership.GetUpdatedAt().AsTime(),
	}
}

// GroupMembershipsListResponse ...
type GroupMembershipsListResponse struct {
	TotalCount       int64                      `json:"totalCount" bson:"total_count"`
	TotalPages       int64                      `json:"totalPages" bson:"total_pages"`
	Page             int64                      `json:"page" bson:"page"`
	Size             int64                      `json:"size" bson:"size"`
	HasMore          bool                       `json:"hasMore" bson:"has_more"`
	GroupMemberships []*GroupMembershipResponse `json:"groupMemberships" bson:"group_memberships"`
}

func GroupMembershipListResponseFromGrpc(listResponse *membershipQueryService.GetGroupMembershipRes) *GroupMembershipsListResponse {
	list := make([]*GroupMembershipResponse, 0, len(listResponse.GetGroupMemberships()))
	for _, gm := range listResponse.GetGroupMemberships() {
		list = append(list, GroupMembershipResponseFromGrpc(gm))
	}
	return &GroupMembershipsListResponse{
		TotalCount:       listResponse.GetTotalCount(),
		TotalPages:       listResponse.GetTotalPages(),
		Page:             listResponse.GetPage(),
		Size:             listResponse.GetSize(),
		HasMore:          listResponse.GetHasMore(),
		GroupMemberships: list,
	}
}
