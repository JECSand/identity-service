package commands

import (
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/gofrs/uuid"
)

// MembershipCommands ...
type MembershipCommands struct {
	CreateMembership CreateMembershipCmdHandler
	UpdateMembership UpdateMembershipCmdHandler
	DeleteMembership DeleteMembershipCmdHandler
}

// NewMembershipCommands ...
func NewMembershipCommands(createMembership CreateMembershipCmdHandler, updateMembership UpdateMembershipCmdHandler, deleteMembership DeleteMembershipCmdHandler) *MembershipCommands {
	return &MembershipCommands{
		CreateMembership: createMembership,
		UpdateMembership: updateMembership,
		DeleteMembership: deleteMembership,
	}
}

// CreateMembershipCommand ...
type CreateMembershipCommand struct {
	ID      uuid.UUID              `json:"id"`
	UserID  uuid.UUID              `json:"userID,omitempty"`
	GroupID uuid.UUID              `json:"groupID,omitempty"`
	Status  enums.MembershipStatus `json:"status,omitempty"`
	Role    enums.Role             `json:"role,omitempty"`
}

// NewCreateMembershipCommand ...
func NewCreateMembershipCommand(id uuid.UUID, userId uuid.UUID, groupId uuid.UUID, status enums.MembershipStatus, role enums.Role) *CreateMembershipCommand {
	return &CreateMembershipCommand{
		ID:      id,
		UserID:  userId,
		GroupID: groupId,
		Status:  status,
		Role:    role,
	}
}

// UpdateMembershipCommand ...
type UpdateMembershipCommand struct {
	ID     uuid.UUID              `json:"id" validate:"required,gte=0,lte=255"`
	Status enums.MembershipStatus `json:"status,omitempty"`
	Role   enums.Role             `json:"role,omitempty"`
}

// NewUpdateMembershipCommand ...
func NewUpdateMembershipCommand(id uuid.UUID, status enums.MembershipStatus, role enums.Role) *UpdateMembershipCommand {
	return &UpdateMembershipCommand{
		ID:     id,
		Status: status,
		Role:   role,
	}
}

// DeleteMembershipCommand ...
type DeleteMembershipCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// NewDeleteMembershipCommand ...
func NewDeleteMembershipCommand(id uuid.UUID) *DeleteMembershipCommand {
	return &DeleteMembershipCommand{ID: id}
}
