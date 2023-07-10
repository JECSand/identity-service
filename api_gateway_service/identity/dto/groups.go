package dto

import (
	groupQueryService "github.com/JECSand/identity-service/query_service/protos/group_query"
	"github.com/gofrs/uuid"
	"time"
)

type CreateGroupDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=255"`
	CreatorID   uuid.UUID `json:"creatorID" validate:"required,gte=0,lte=5000"`
	Active      bool      `json:"active"`
}

type CreateGroupResponseDTO struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type UpdateGroupDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=255"`
}

// GroupResponse ...
type GroupResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatorID   string    `json:"creatorID,omitempty"`
	Active      bool      `json:"active,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

func GroupResponseFromGrpc(group *groupQueryService.Group) *GroupResponse {
	return &GroupResponse{
		ID:          group.GetID(),
		Name:        group.GetName(),
		Description: group.GetDescription(),
		CreatorID:   group.GetCreatorID(),
		Active:      group.GetActive(),
		CreatedAt:   group.GetCreatedAt().AsTime(),
		UpdatedAt:   group.GetUpdatedAt().AsTime(),
	}
}

// GroupsListResponse ...
type GroupsListResponse struct {
	TotalCount int64            `json:"totalCount" bson:"total_count"`
	TotalPages int64            `json:"totalPages" bson:"total_pages"`
	Page       int64            `json:"page" bson:"page"`
	Size       int64            `json:"size" bson:"size"`
	HasMore    bool             `json:"hasMore" bson:"has_more"`
	Groups     []*GroupResponse `json:"groups" bson:"groups"`
}

func GroupsListResponseFromGrpc(listResponse *groupQueryService.SearchGroupRes) *GroupsListResponse {
	list := make([]*GroupResponse, 0, len(listResponse.GetGroups()))
	for _, group := range listResponse.GetGroups() {
		list = append(list, GroupResponseFromGrpc(group))
	}
	return &GroupsListResponse{
		TotalCount: listResponse.GetTotalCount(),
		TotalPages: listResponse.GetTotalPages(),
		Page:       listResponse.GetPage(),
		Size:       listResponse.GetSize(),
		HasMore:    listResponse.GetHasMore(),
		Groups:     list,
	}
}
