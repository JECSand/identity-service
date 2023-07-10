package commands

import (
	"github.com/gofrs/uuid"
)

// GroupCommands ...
type GroupCommands struct {
	CreateGroup CreateGroupCmdHandler
	UpdateGroup UpdateGroupCmdHandler
	DeleteGroup DeleteGroupCmdHandler
}

// NewGroupCommands ...
func NewGroupCommands(createUser CreateGroupCmdHandler, updateUser UpdateGroupCmdHandler, deleteUser DeleteGroupCmdHandler) *GroupCommands {
	return &GroupCommands{
		CreateGroup: createUser,
		UpdateGroup: updateUser,
		DeleteGroup: deleteUser,
	}
}

// CreateGroupCommand ...
type CreateGroupCommand struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	CreatorID   uuid.UUID `json:"creatorID" validate:"required"`
	Active      bool      `json:"active"`
}

// NewCreateGroupCommand ...
func NewCreateGroupCommand(id uuid.UUID, name string, description string, creatorId uuid.UUID, active bool) *CreateGroupCommand {
	return &CreateGroupCommand{
		ID:          id,
		Name:        name,
		Description: description,
		CreatorID:   creatorId,
		Active:      active,
	}
}

// UpdateGroupCommand ...
type UpdateGroupCommand struct {
	ID          uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// NewUpdateGroupCommand ...
func NewUpdateGroupCommand(id uuid.UUID, name string, description string) *UpdateGroupCommand {
	return &UpdateGroupCommand{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

// DeleteGroupCommand ...
type DeleteGroupCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// NewDeleteGroupCommand ...
func NewDeleteGroupCommand(id uuid.UUID) *DeleteGroupCommand {
	return &DeleteGroupCommand{ID: id}
}
