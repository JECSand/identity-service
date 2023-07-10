package models

import (
	"github.com/gofrs/uuid"
	"time"
)

// Blacklist is a root struct
type Blacklist struct {
	ID          uuid.UUID `json:"id"`
	AccessToken string    `json:"accessToken,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
