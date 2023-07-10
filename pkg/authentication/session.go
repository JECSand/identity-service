package authentication

import (
	"errors"
	"fmt"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Session stores the structured data from a session token for use
type Session struct {
	UserId     string
	RootAdmin  bool
	Type       enums.SessionType
	Expiration int64
	Cfg        *Config
}

// newSession constructions a new Token
func newSession(userId string, root bool, tokenType enums.SessionType, cfg *Config) *Session {
	return &Session{
		UserId:     userId,
		RootAdmin:  root,
		Type:       tokenType,
		Expiration: 0,
		Cfg:        cfg,
	}
}

// setExpiration returns the unix time for token expiration
func (t *Session) setExpiration() {
	var duration int
	switch {
	case t.Type == enums.USER:
		duration = t.Cfg.SessionDuration
	case t.Type == enums.INTEGRATION:
		duration = t.Cfg.IntegrationDuration
	default:
		duration = 1
	}
	t.Expiration = time.Now().Add(time.Hour * time.Duration(duration)).Unix()
}

// NewToken is used to create an encrypted token string from an auth Session
func (t *Session) NewToken() (string, error) {
	if t.UserId == "" {
		return "", errors.New("missing required token claims")
	}
	t.setExpiration()
	if t.Expiration == 0 {
		return "", errors.New("new token must have a expiration time greater than 0")
	}
	var MySigningKey = []byte(t.Cfg.Secret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = t.UserId
	claims["root"] = t.RootAdmin
	claims["token_type"] = t.Type.Stringify()
	claims["exp"] = t.Expiration
	return token.SignedString(MySigningKey)
}

// decryptToken a Session from an encrypted token string
func decryptToken(tokenStr string, secret string) (*Session, error) {
	var session Session
	if tokenStr == "" {
		return &session, errors.New("unauthorized")
	}
	var MySigningKey = []byte(secret)
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return MySigningKey, nil
	})
	if err != nil {
		return &session, err
	}
	if parsedToken.Valid {
		tokenClaims := parsedToken.Claims.(jwt.MapClaims)
		session.UserId = tokenClaims["id"].(string)
		session.RootAdmin = tokenClaims["root"].(bool)
		session.Type = enums.SessionTypeFromString(tokenClaims["token_type"].(string))
		return &session, nil
	}
	return &session, errors.New("invalid token")
}
