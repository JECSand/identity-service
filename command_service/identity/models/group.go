package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Group struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatorID   uuid.UUID `json:"creatorID,omitempty"`
	Active      bool      `json:"active,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
