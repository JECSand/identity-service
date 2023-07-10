package queries

// AuthQueries ...
type AuthQueries struct {
	CheckTokenBlacklist CheckTokenBlacklistHandler
}

// NewAuthQueries ...
func NewAuthQueries(checkBlacklist CheckTokenBlacklistHandler) *AuthQueries {
	return &AuthQueries{
		CheckTokenBlacklist: checkBlacklist,
	}
}

// CheckTokenBlacklistQuery ...
type CheckTokenBlacklistQuery struct {
	AccessToken string `json:"accessToken" validate:"required,gte=0,lte=255"`
}

// NewCheckTokenBlacklistQuery ...
func NewCheckTokenBlacklistQuery(accessToken string) *CheckTokenBlacklistQuery {
	return &CheckTokenBlacklistQuery{
		AccessToken: accessToken,
	}
}
