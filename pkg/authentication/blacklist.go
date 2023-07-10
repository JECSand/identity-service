package authentication

import "time"

type BlacklistService interface {
	BlacklistAuthToken(accessToken string) error
	CheckTokenBlacklist(accessToken string) bool
}

// Blacklist is a root struct that is used to store the json encoded data for/from a mongodb blacklist doc.
type Blacklist struct {
	Id          string    `json:"id,omitempty"`
	AccessToken string    `json:"accessToken,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}
