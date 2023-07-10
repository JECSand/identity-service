package events

import (
	"time"
)

type AuthEvents struct {
	BlacklistToken BlacklistTokenEventHandler
	UpdatePassword UpdatePasswordEventHandler
}

func NewAuthEvents(
	blacklistToken BlacklistTokenEventHandler,
	updatePassword UpdatePasswordEventHandler,
) *AuthEvents {
	return &AuthEvents{
		BlacklistToken: blacklistToken,
		UpdatePassword: updatePassword,
	}
}

type BlacklistTokenEvent struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	AccessToken string    `json:"accessToken,omitempty" bson:"access_token,omitempty" validate:"required,min=3,max=250"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NewBlacklistTokenEvent(id string, accessToken string, createdAt time.Time, updatedAt time.Time) *BlacklistTokenEvent {
	return &BlacklistTokenEvent{
		ID:          id,
		AccessToken: accessToken,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type UpdatePasswordEvent struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	NewPassword string    `json:"newPassword,omitempty" bson:"new_password,omitempty" validate:"required,min=3,max=250"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NewUpdatePasswordEvent(id string, newPassword string, up time.Time) *UpdatePasswordEvent {
	return &UpdatePasswordEvent{
		ID:          id,
		NewPassword: newPassword,
		UpdatedAt:   up,
	}
}
