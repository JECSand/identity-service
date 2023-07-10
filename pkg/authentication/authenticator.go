package authentication

import (
	"context"
	"errors"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// Config settings for auth
type Config struct {
	SessionDuration     int    `mapstructure:"sessionDuration"`     // 1
	IntegrationDuration int    `mapstructure:"integrationDuration"` // 4380
	Secret              string `mapstructure:"secret"`              // 4380
}

func NewAuthConfig(uDur int, iDur int, secret string) *Config {
	return &Config{
		SessionDuration:     uDur,
		IntegrationDuration: iDur,
		Secret:              secret,
	}
}

type Authenticator interface {
	NewSession(userId string, root bool, tokenType enums.SessionType) *Session
	GetTokenSession(accessToken string) (*Session, error)
	AuthorizeGRPC(ctx context.Context, method string) (*Session, error)
	AuthorizeREST(req *http.Request, method string) (*Session, error)
}

// authenticator
type authenticator struct {
	log         logging.Logger
	accessRules map[string]enums.Role
	cfg         *Config
}

// NewAuthenticator constructs a new authenticator
func NewAuthenticator(log logging.Logger, aRules map[string]enums.Role, cfg *Config) *authenticator {
	return &authenticator{
		log:         log,
		accessRules: aRules,
		cfg:         cfg,
	}
}

// authorize a gRPC request
func (i *authenticator) authorize(accessToken string, accessRule enums.Role) (*Session, error) {
	var authorized bool
	session, err := i.GetTokenSession(accessToken)
	if err != nil {
		return session, err
	}
	switch accessRule {
	case enums.ROOT:
		if session.RootAdmin {
			authorized = true
		}
	case enums.MEMBER:
		authorized = true
	}
	if err != nil {
		return session, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	if authorized {
		return session, nil
	}
	return session, status.Error(codes.PermissionDenied, "no permission to access this controller")
}

// GetTokenSession validates & decrypts a JWT token, then returns the Session
func (i *authenticator) GetTokenSession(accessToken string) (*Session, error) {
	tokenSession, err := decryptToken(accessToken, i.cfg.Secret)
	if err != nil {
		return nil, err
	}
	return tokenSession, nil
}

// NewSession validates & decrypts a JWT token, then returns the Session
func (i *authenticator) NewSession(userId string, root bool, tokenType enums.SessionType) *Session {
	return newSession(userId, root, tokenType, i.cfg)
}

// AuthorizeREST a REST request
func (i *authenticator) AuthorizeREST(req *http.Request, method string) (*Session, error) {
	accessRule, ok := i.accessRules[method]
	if !ok {
		return nil, nil // unprotected endpoint
	}
	accessToken := req.Header.Get("Authorization")
	if accessToken == "" {
		return nil, errors.New("unauthorized")
	}
	return i.authorize(accessToken, accessRule)
}

// AuthorizeGRPC a gRPC request
func (i *authenticator) AuthorizeGRPC(ctx context.Context, method string) (*Session, error) {
	accessRule, ok := i.accessRules[method]
	if !ok {
		return nil, nil // unprotected endpoint
	}
	accessToken, err := GetTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return i.authorize(accessToken, accessRule)

}
