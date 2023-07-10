package dto

import (
	authQueryService "github.com/JECSand/identity-service/query_service/protos/auth_query"
	"github.com/gofrs/uuid"
	"time"
)

type RegisterDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username,omitempty" validate:"required,gte=0,lte=255"`
	Root     bool      `json:"root,omitempty" validate:"required,gte=0,lte=255"`
	Email    string    `json:"email" validate:"required,gte=0,lte=255"`
	Password string    `json:"password" validate:"required,gte=0,lte=5000"`
}

// AuthUserResponse ...
type AuthUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email,omitempty"`
	Username  string    `json:"username,omitempty"`
	Root      bool      `json:"root,omitempty"`
	Active    bool      `json:"active,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type AuthenticateDTO struct {
	Email    string `json:"email" validate:"required,gte=0,lte=255"`
	Password string `json:"password" validate:"required,gte=0,lte=5000"`
}

type AuthenticateResponse struct {
	User   *AuthUserResponse `json:"user" validate:"required,gte=0,lte=255"`
	Status int64             `json:"status" validate:"required,gte=0,lte=255"`
}

func AuthUserResponseFromGrpc(aUser *authQueryService.User) *AuthUserResponse {
	return &AuthUserResponse{
		ID:        aUser.GetID(),
		Email:     aUser.GetEmail(),
		Username:  aUser.GetUsername(),
		Root:      aUser.GetRoot(),
		Active:    aUser.GetActive(),
		CreatedAt: aUser.GetCreatedAt().AsTime(),
		UpdatedAt: aUser.GetUpdatedAt().AsTime(),
	}
}

func AuthenticateResponseResponseFromGrpc(auth *authQueryService.AuthenticateRes) *AuthenticateResponse {
	return &AuthenticateResponse{
		User:   AuthUserResponseFromGrpc(auth.GetUser()),
		Status: auth.GetStatus(),
	}
}

type ValidateResponse struct {
	User   *AuthUserResponse `json:"user" validate:"required,gte=0,lte=255"`
	Status int64             `json:"status" validate:"required,gte=0,lte=255"`
}

func ValidateResponseResponseFromGrpc(auth *authQueryService.ValidateRes) *ValidateResponse {
	return &ValidateResponse{
		User:   AuthUserResponseFromGrpc(auth.GetUser()),
		Status: auth.GetStatus(),
	}
}

type BlacklistTokenDTO struct {
	ID          uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
	AccessToken string    `json:"accessToken" validate:"required,gte=0,lte=255"`
}

type UpdatePasswordDTO struct {
	ID              uuid.UUID `json:"id" validate:"required,gte=0,lte=255"`
	CurrentPassword string    `json:"currentPassword" validate:"required,gte=0,lte=255"`
	NewPassword     string    `json:"newPassword" validate:"required,gte=0,lte=255"`
}

type ErrorDTO struct {
	Message string `json:"message" validate:"required,gte=0,lte=255"`
}
