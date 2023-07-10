package commands

import (
	"github.com/gofrs/uuid"
)

// AuthCommands ...
type AuthCommands struct {
	BlacklistToken BlacklistTokenCmdHandler
	UpdatePassword PasswordUpdateCmdHandler
}

// NewAuthCommands ...
func NewAuthCommands(blacklistToken BlacklistTokenCmdHandler, passwordUpdate PasswordUpdateCmdHandler) *AuthCommands {
	return &AuthCommands{
		BlacklistToken: blacklistToken,
		UpdatePassword: passwordUpdate,
	}
}

// BlacklistTokenCommand ...
type BlacklistTokenCommand struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	AccessToken string    `json:"accessToken" validate:"required,gte=0,lte=255"`
}

// NewBlacklistTokenCommand ...
func NewBlacklistTokenCommand(id uuid.UUID, accessToken string) *BlacklistTokenCommand {
	return &BlacklistTokenCommand{
		ID:          id,
		AccessToken: accessToken,
	}
}

// PasswordUpdateCommand ...
type PasswordUpdateCommand struct {
	ID              uuid.UUID `json:"id" validate:"required"`
	CurrentPassword string    `json:"currentPassword" validate:"required"`
	NewPassword     string    `json:"newPassword" validate:"required"`
}

// NewUpdatePasswordCommand ...
func NewUpdatePasswordCommand(id uuid.UUID, currentPassword string, newPassword string) *PasswordUpdateCommand {
	return &PasswordUpdateCommand{
		ID:              id,
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}
}
