package services

import (
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/cache"
	"github.com/JECSand/identity-service/query_service/identity/data"
	"github.com/JECSand/identity-service/query_service/identity/events"
	"github.com/JECSand/identity-service/query_service/identity/queries"
)

type AuthService struct {
	Events  *events.AuthEvents
	Queries *queries.AuthQueries
}

func NewAuthService(
	log logging.Logger,
	cfg *config.Config,
	mongoDB data.Database,
	redisCache cache.Cache,
) *AuthService {
	blacklistTokenHandler := events.NewBlacklistTokenEventHandler(log, cfg, mongoDB, redisCache)
	updatePasswordEventHandler := events.NewUpdatePasswordEventHandler(log, cfg, mongoDB, redisCache)
	authenticateHandler := queries.NewAuthenticateHandler(log, cfg, mongoDB, redisCache)
	validateHandler := queries.NewValidateHandler(log, cfg, mongoDB, redisCache)
	userEvents := events.NewAuthEvents(blacklistTokenHandler, updatePasswordEventHandler)
	userQueries := queries.NewAuthQueries(authenticateHandler, validateHandler)
	return &AuthService{
		Events:  userEvents,
		Queries: userQueries,
	}
}
