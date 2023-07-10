package models

import (
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/gofrs/uuid"
	"time"
)

type Membership struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"userID,omitempty"`
	GroupID   uuid.UUID              `json:"groupID,omitempty"`
	Status    enums.MembershipStatus `json:"status,omitempty"`
	Role      enums.Role             `json:"role,omitempty"`
	CreatedAt time.Time              `json:"createdAt,omitempty"`
	UpdatedAt time.Time              `json:"updatedAt,omitempty"`
}

type UserMembership struct {
	ID           uuid.UUID              `json:"id"`
	GroupID      uuid.UUID              `json:"groupID,omitempty"`
	UserID       uuid.UUID              `json:"userID,omitempty"`
	MembershipID uuid.UUID              `json:"membershipID,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Username     string                 `json:"username,omitempty"`
	Status       enums.MembershipStatus `json:"status,omitempty"`
	Role         enums.Role             `json:"role,omitempty"`
	CreatedAt    time.Time              `json:"createdAt,omitempty"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty"`
}

type GroupMembership struct {
	ID           uuid.UUID              `json:"id"`
	UserID       uuid.UUID              `json:"userID,omitempty"`
	GroupID      uuid.UUID              `json:"groupID,omitempty"`
	MembershipID uuid.UUID              `json:"membershipID,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Status       enums.MembershipStatus `json:"status,omitempty"`
	Role         enums.Role             `json:"role,omitempty"`
	Creator      bool                   `json:"creator,omitempty"`
	CreatedAt    time.Time              `json:"createdAt,omitempty"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty"`
}
