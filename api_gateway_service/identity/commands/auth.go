package commands

import (
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
)

type AuthCommands struct {
	BlacklistToken BlacklistTokenCmdHandler
	UpdatePassword UpdatePasswordCmdHandler
}

func NewAuthCommands(blacklistToken BlacklistTokenCmdHandler, updatePass UpdatePasswordCmdHandler) *AuthCommands {
	return &AuthCommands{
		BlacklistToken: blacklistToken,
		UpdatePassword: updatePass,
	}
}

// BlacklistTokenCommand ...
type BlacklistTokenCommand struct {
	BlacklistDto *dto.BlacklistTokenDTO
}

func NewBlacklistTokenCommand(blacklistDto *dto.BlacklistTokenDTO) *BlacklistTokenCommand {
	return &BlacklistTokenCommand{BlacklistDto: blacklistDto}
}

// UpdatePasswordCommand ...
type UpdatePasswordCommand struct {
	UpdateDto *dto.UpdatePasswordDTO
}

func NewUpdatePasswordCommand(updateDto *dto.UpdatePasswordDTO) *UpdatePasswordCommand {
	return &UpdatePasswordCommand{UpdateDto: updateDto}
}
