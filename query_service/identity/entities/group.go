package entities

import (
	"github.com/JECSand/identity-service/pkg/utilities"
	queryService "github.com/JECSand/identity-service/query_service/protos/group_query"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Group struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"email,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	CreatorID   string    `json:"creatorID,omitempty" bson:"creator_id,omitempty" validate:"required"`
	Active      bool      `json:"active,omitempty" bson:"active,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// GetID returns the unique identifier of the Group
func (u *Group) GetID() string {
	return u.ID
}

// GroupsList response with pagination
type GroupsList struct {
	TotalCount int64    `json:"totalCount" bson:"total_count"`
	TotalPages int64    `json:"totalPages" bson:"total_pages"`
	Page       int64    `json:"page" bson:"page"`
	Size       int64    `json:"size" bson:"size"`
	HasMore    bool     `json:"hasMore" bson:"has_more"`
	Groups     []*Group `json:"groups" bson:"groups"`
}

func NewGroupListWithPagination(groups []*Group, count int64, pagination *utilities.Pagination) *GroupsList {
	return &GroupsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Groups:     groups,
	}
}

func GroupToGrpcMessage(group *Group) *queryService.Group {
	return &queryService.Group{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		CreatorID:   group.CreatorID,
		Active:      group.Active,
		CreatedAt:   timestamppb.New(group.CreatedAt),
		UpdatedAt:   timestamppb.New(group.UpdatedAt),
	}
}

func GroupListToGrpc(groups *GroupsList) *queryService.SearchGroupRes {
	list := make([]*queryService.Group, 0, len(groups.Groups))
	for _, group := range groups.Groups {
		list = append(list, GroupToGrpcMessage(group))
	}
	return &queryService.SearchGroupRes{
		TotalCount: groups.TotalCount,
		TotalPages: groups.TotalPages,
		Page:       groups.Page,
		Size:       groups.Size,
		HasMore:    groups.HasMore,
		Groups:     list,
	}
}
