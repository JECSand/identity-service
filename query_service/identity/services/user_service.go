package services

import (
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/cache"
	"github.com/JECSand/identity-service/query_service/identity/data"
	"github.com/JECSand/identity-service/query_service/identity/events"
	"github.com/JECSand/identity-service/query_service/identity/queries"
)

type UserService struct {
	Events  *events.UserEvents
	Queries *queries.UserQueries
}

func NewUserService(
	log logging.Logger,
	cfg *config.Config,
	mongoDB data.Database,
	redisCache cache.Cache,
) *UserService {
	createUserHandler := events.NewCreateUserEventHandler(log, cfg, mongoDB, redisCache)
	deleteUserEventHandler := events.NewDeleteUserEventHandler(log, cfg, mongoDB, redisCache)
	updateUserEventHandler := events.NewUpdateUserEventHandler(log, cfg, mongoDB, redisCache)
	getUserByIdHandler := queries.NewGetUserByIdHandler(log, cfg, mongoDB, redisCache)
	searchUserHandler := queries.NewSearchUserHandler(log, cfg, mongoDB, redisCache)
	userEvents := events.NewUserEvents(createUserHandler, updateUserEventHandler, deleteUserEventHandler)
	userQueries := queries.NewUserQueries(getUserByIdHandler, searchUserHandler)
	return &UserService{
		Events:  userEvents,
		Queries: userQueries,
	}
}
