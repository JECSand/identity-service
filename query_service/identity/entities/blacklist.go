package entities

import (
	"time"
)

// Blacklist is a root struct that is used to store the json encoded data for/from a mongodb blacklist doc.
type Blacklist struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	AccessToken string    `json:"accessToken,omitempty" bson:"access_token,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// GetID returns the unique identifier of the Blacklist
func (b *Blacklist) GetID() string {
	return b.ID
}
