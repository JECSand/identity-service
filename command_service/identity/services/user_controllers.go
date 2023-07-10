package services

import (
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	"github.com/JECSand/identity-service/command_service/identity/queries"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
)

// UserService ...
type UserService struct {
	Commands *commands.UserCommands
	Queries  *queries.UserQueries
}

// NewUserService ...
func NewUserService(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *UserService {
	updateUserHandler := commands.NewUpdateUserHandler(log, cfg, pgRepo, kafkaProducer)
	createUserHandler := commands.NewCreateUserHandler(log, cfg, pgRepo, kafkaProducer)
	deleteUserHandler := commands.NewDeleteUserHandler(log, cfg, pgRepo, kafkaProducer)
	getUserByIdHandler := queries.NewGetUserByIdHandler(log, cfg, pgRepo)
	countUsersHandler := queries.NewCountUsersHandler(log, cfg, pgRepo)
	userCommands := commands.NewUserCommands(createUserHandler, updateUserHandler, deleteUserHandler)
	userQueries := queries.NewUserQueries(getUserByIdHandler, countUsersHandler)
	return &UserService{
		Commands: userCommands,
		Queries:  userQueries,
	}
}
