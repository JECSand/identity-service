package commands

import (
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/gofrs/uuid"
)

type MembershipCommands struct {
	CreateMembership CreateMembershipCmdHandler
	UpdateMembership UpdateMembershipCmdHandler
	DeleteMembership DeleteMembershipCmdHandler
}

func NewMembershipCommands(create CreateMembershipCmdHandler, update UpdateMembershipCmdHandler, delete DeleteMembershipCmdHandler) *MembershipCommands {
	return &MembershipCommands{
		CreateMembership: create,
		UpdateMembership: update,
		DeleteMembership: delete,
	}
}

// CreateMembershipCommand ...
type CreateMembershipCommand struct {
	CreateDto *dto.CreateMembershipDTO
}

func NewCreateMembershipCommand(createDto *dto.CreateMembershipDTO) *CreateMembershipCommand {
	return &CreateMembershipCommand{CreateDto: createDto}
}

// UpdateMembershipCommand ...
type UpdateMembershipCommand struct {
	UpdateDto *dto.UpdateMembershipDTO
}

func NewUpdateMembershipCommand(updateDto *dto.UpdateMembershipDTO) *UpdateMembershipCommand {
	return &UpdateMembershipCommand{UpdateDto: updateDto}
}

// DeleteMembershipCommand ...
type DeleteMembershipCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

func NewDeleteMembershipCommand(membershipID uuid.UUID) *DeleteMembershipCommand {
	return &DeleteMembershipCommand{ID: membershipID}
}
