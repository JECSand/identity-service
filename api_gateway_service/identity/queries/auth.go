package queries

import "github.com/JECSand/identity-service/pkg/enums"

type AuthQueries struct {
	Authenticate AuthenticateHandler
	Validate     ValidateHandler
}

func NewAuthQueries(authenticate AuthenticateHandler, validate ValidateHandler) *AuthQueries {
	return &AuthQueries{
		Authenticate: authenticate,
		Validate:     validate,
	}
}

type AuthenticateQuery struct {
	Email    string `json:"email" validate:"required,gte=0,lte=255"`
	Password string `json:"password" validate:"required,gte=0,lte=255"`
}

func NewAuthenticateQuery(email string, password string) *AuthenticateQuery {
	return &AuthenticateQuery{
		Email:    email,
		Password: password,
	}
}

type ValidateQuery struct {
	UserID         string               `json:"userID validate:required,gte=0,lte=255"`
	AccessToken    string               `json:"accessToken validate:required,gte=0,lte=255"`
	ValidationType enums.ValidationType `json:"validationType validate:required,gte=0,lte=255"`
}

func NewValidateQuery(userID string, accessToken string, valType enums.ValidationType) *ValidateQuery {
	return &ValidateQuery{
		UserID:         userID,
		AccessToken:    accessToken,
		ValidationType: valType,
	}
}
