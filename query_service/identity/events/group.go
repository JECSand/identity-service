package events

import (
	"github.com/gofrs/uuid"
	"time"
)

type GroupEvents struct {
	CreateGroup CreateGroupEventHandler
	UpdateGroup UpdateGroupEventHandler
	DeleteGroup DeleteGroupEventHandler
}

func NewGroupEvents(
	createGroup CreateGroupEventHandler,
	updateGroup UpdateGroupEventHandler,
	deleteGroup DeleteGroupEventHandler,
) *GroupEvents {
	return &GroupEvents{
		CreateGroup: createGroup,
		UpdateGroup: updateGroup,
		DeleteGroup: deleteGroup,
	}
}

type CreateGroupEvent struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	CreatorID   string    `json:"creatorID,omitempty" bson:"creator_id,omitempty" validate:"required"`
	Active      bool      `json:"active,omitempty" bson:"active,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NewCreateGroupEvent(id string, name string, description string, creatorID string, active bool, createdAt time.Time, updatedAt time.Time) *CreateGroupEvent {
	return &CreateGroupEvent{
		ID:          id,
		Name:        name,
		Description: description,
		CreatorID:   creatorID,
		Active:      active,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type UpdateGroupEvent struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NewUpdateGroupEvent(id string, name string, description string, updatedAt time.Time) *UpdateGroupEvent {
	return &UpdateGroupEvent{
		ID:          id,
		Name:        name,
		Description: description,
		UpdatedAt:   updatedAt,
	}
}

type DeleteGroupEvent struct {
	ID uuid.UUID `json:"id" bson:"_id,omitempty"`
}

func NewDeleteGroupEvent(id uuid.UUID) *DeleteGroupEvent {
	return &DeleteGroupEvent{ID: id}
}
