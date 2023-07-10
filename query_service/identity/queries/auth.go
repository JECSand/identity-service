package queries

import (
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/gofrs/uuid"
)

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
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
}

func NewAuthenticateQuery(email string, password string) *AuthenticateQuery {
	return &AuthenticateQuery{
		Email:    email,
		Password: password,
	}
}

type ValidateQuery struct {
	UserID         uuid.UUID            `json:"userID" bson:"_id,omitempty"`
	AccessToken    string               `json:"accessToken" bson:"accessToken,omitempty"`
	ValidationType enums.ValidationType `json:"validationType" bson:"accessToken,omitempty"`
}

func NewValidateQuery(userID uuid.UUID, accessToken string, valType enums.ValidationType) *ValidateQuery {
	return &ValidateQuery{
		UserID:         userID,
		AccessToken:    accessToken,
		ValidationType: valType,
	}
}
