package commands

import (
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/gofrs/uuid"
)

type GroupCommands struct {
	CreateGroup CreateGroupCmdHandler
	UpdateGroup UpdateGroupCmdHandler
	DeleteGroup DeleteGroupCmdHandler
}

func NewGroupCommands(create CreateGroupCmdHandler, update UpdateGroupCmdHandler, delete DeleteGroupCmdHandler) *GroupCommands {
	return &GroupCommands{
		CreateGroup: create,
		UpdateGroup: update,
		DeleteGroup: delete,
	}
}

// CreateGroupCommand ...
type CreateGroupCommand struct {
	CreateDto *dto.CreateGroupDTO
}

func NewCreateGroupCommand(createDto *dto.CreateGroupDTO) *CreateGroupCommand {
	return &CreateGroupCommand{CreateDto: createDto}
}

// UpdateGroupCommand ...
type UpdateGroupCommand struct {
	UpdateDto *dto.UpdateGroupDTO
}

func NewUpdateGroupCommand(updateDto *dto.UpdateGroupDTO) *UpdateGroupCommand {
	return &UpdateGroupCommand{UpdateDto: updateDto}
}

// DeleteGroupCommand ...
type DeleteGroupCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

func NewDeleteGroupCommand(groupID uuid.UUID) *DeleteGroupCommand {
	return &DeleteGroupCommand{ID: groupID}
}
