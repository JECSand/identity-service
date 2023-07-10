package commands

import (
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/gofrs/uuid"
)

type UserCommands struct {
	CreateUser CreateUserCmdHandler
	UpdateUser UpdateUserCmdHandler
	DeleteUser DeleteUserCmdHandler
}

func NewUserCommands(create CreateUserCmdHandler, update UpdateUserCmdHandler, delete DeleteUserCmdHandler) *UserCommands {
	return &UserCommands{
		CreateUser: create,
		UpdateUser: update,
		DeleteUser: delete,
	}
}

// CreateUserCommand ...
type CreateUserCommand struct {
	CreateDto *dto.CreateUserDTO
}

func NewCreateUserCommand(createDto *dto.CreateUserDTO) *CreateUserCommand {
	return &CreateUserCommand{CreateDto: createDto}
}

// UpdateUserCommand ...
type UpdateUserCommand struct {
	UpdateDto *dto.UpdateUserDTO
}

func NewUpdateUserCommand(updateDto *dto.UpdateUserDTO) *UpdateUserCommand {
	return &UpdateUserCommand{UpdateDto: updateDto}
}

// DeleteUserCommand ...
type DeleteUserCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

func NewDeleteUserCommand(userID uuid.UUID) *DeleteUserCommand {
	return &DeleteUserCommand{ID: userID}
}
