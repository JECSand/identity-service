package events

import (
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/gofrs/uuid"
	"time"
)

type MembershipEvents struct {
	CreateMembership CreateMembershipEventHandler
	UpdateMembership UpdateMembershipEventHandler
	DeleteMembership DeleteMembershipEventHandler
}

func NewMembershipEvents(
	createMembership CreateMembershipEventHandler,
	updateMembership UpdateMembershipEventHandler,
	deleteMembership DeleteMembershipEventHandler,
) *MembershipEvents {
	return &MembershipEvents{
		CreateMembership: createMembership,
		UpdateMembership: updateMembership,
		DeleteMembership: deleteMembership,
	}
}

type CreatedUserMembership struct {
	ID           string                 `json:"id" bson:"_id,omitempty" validate:"required"`
	GroupID      string                 `json:"groupID,omitempty" bson:"group_id,omitempty" validate:"required,min=3,max=250"`
	UserID       string                 `json:"userID,omitempty" bson:"user_id,omitempty" validate:"required,min=3,max=500"`
	MembershipID string                 `json:"membershipID,omitempty" bson:"membership_id,omitempty" validate:"required,min=3,max=500"`
	Email        string                 `json:"email,omitempty" bson:"email,omitempty" validate:"required,min=3,max=500"`
	Username     string                 `json:"username,omitempty" bson:"username,omitempty" validate:"required,min=3,max=500"`
	Status       enums.MembershipStatus `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
	Role         enums.Role             `json:"role,omitempty" bson:"role,omitempty" validate:"required"`
	CreatedAt    time.Time              `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NewCreatedUserMembership(
	id string,
	groupID string,
	userID string,
	membershipID string,
	email string,
	username string,
	status enums.MembershipStatus,
	role enums.Role,
	createdAt time.Time,
	updatedAt time.Time,
) *CreatedUserMembership {
	return &CreatedUserMembership{
		ID:           id,
		GroupID:      groupID,
		UserID:       userID,
		MembershipID: membershipID,
		Email:        email,
		Username:     username,
		Status:       status,
		Role:         role,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

type CreatedGroupMembership struct {
	ID           string                 `json:"id" bson:"_id,omitempty" validate:"required"`
	UserID       string                 `json:"userID,omitempty" bson:"user_id,omitempty" validate:"required,min=3,max=250"`
	GroupID      string                 `json:"groupID,omitempty" bson:"group_id,omitempty" validate:"required,min=3,max=500"`
	MembershipID string                 `json:"membershipID,omitempty" bson:"membership_id,omitempty" validate:"required,min=3,max=500"`
	Name         string                 `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=500"`
	Description  string                 `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Status       enums.MembershipStatus `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
	Role         enums.Role             `json:"role,omitempty" bson:"role,omitempty" validate:"required"`
	Creator      bool                   `json:"creator,omitempty" bson:"creator,omitempty"`
	CreatedAt    time.Time              `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NewCreatedGroupMembership(
	id string,
	groupID string,
	userID string,
	membershipID string,
	name string,
	description string,
	status enums.MembershipStatus,
	role enums.Role,
	creator bool,
	createdAt time.Time,
	updatedAt time.Time,
) *CreatedGroupMembership {
	return &CreatedGroupMembership{
		ID:           id,
		GroupID:      groupID,
		UserID:       userID,
		MembershipID: membershipID,
		Name:         name,
		Description:  description,
		Status:       status,
		Role:         role,
		Creator:      creator,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

type CreatedMembership struct {
	ID        string                 `json:"id" bson:"_id,omitempty" validate:"required"`
	UserID    string                 `json:"userID,omitempty" bson:"user_id,omitempty" validate:"required,min=3,max=250"`
	GroupID   string                 `json:"groupID,omitempty" bson:"group_id,omitempty" validate:"required,min=3,max=500"`
	Status    enums.MembershipStatus `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
	Role      enums.Role             `json:"role,omitempty" bson:"role,omitempty" validate:"required"`
	CreatedAt time.Time              `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time              `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NewCreatedMembership(
	id string,
	userID string,
	groupID string,
	status enums.MembershipStatus,
	role enums.Role,
	createdAt time.Time,
	updatedAt time.Time,
) *CreatedMembership {
	return &CreatedMembership{
		ID:        id,
		UserID:    userID,
		GroupID:   groupID,
		Status:    status,
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type CreateMembershipEvent struct {
	Membership      *CreatedMembership
	UserMembership  *CreatedUserMembership
	GroupMembership *CreatedGroupMembership
}

func NewCreateMembershipEvent(membership *CreatedMembership, userMembership *CreatedUserMembership, groupMembership *CreatedGroupMembership) *CreateMembershipEvent {
	return &CreateMembershipEvent{
		Membership:      membership,
		UserMembership:  userMembership,
		GroupMembership: groupMembership,
	}
}

type UpdateMembershipEvent struct {
	ID        string                 `json:"id" bson:"_id,omitempty" validate:"required"`
	Status    enums.MembershipStatus `json:"status,omitempty" bson:"status,omitempty"`
	Role      enums.Role             `json:"role,omitempty" bson:"role,omitempty"`
	UpdatedAt time.Time              `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NewUpdateMembershipEvent(id string, status enums.MembershipStatus, role enums.Role, updatedAt time.Time) *UpdateMembershipEvent {
	return &UpdateMembershipEvent{
		ID:        id,
		Status:    status,
		Role:      role,
		UpdatedAt: updatedAt,
	}
}

type DeleteMembershipEvent struct {
	ID uuid.UUID `json:"id" bson:"_id,omitempty" validate:"required"`
}

func NewDeleteMembershipEvent(id uuid.UUID) *DeleteMembershipEvent {
	return &DeleteMembershipEvent{ID: id}
}
