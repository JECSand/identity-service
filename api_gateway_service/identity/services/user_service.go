package services

import (
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/commands"
	"github.com/JECSand/identity-service/api_gateway_service/identity/queries"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	queryService "github.com/JECSand/identity-service/query_service/protos/user_query"
)

type UserService struct {
	Commands *commands.UserCommands
	Queries  *queries.UserQueries
}

func NewUserService(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer, rsClient queryService.QueryServiceClient) *UserService {
	createUserHandler := commands.NewCreateUserHandler(log, cfg, kafkaProducer)
	updateUserHandler := commands.NewUpdateUserHandler(log, cfg, kafkaProducer)
	deleteUserHandler := commands.NewDeleteUserHandler(log, cfg, kafkaProducer)
	getUserByIdHandler := queries.NewGetUserByIdHandler(log, cfg, rsClient)
	searchUserHandler := queries.NewSearchUserHandler(log, cfg, rsClient)
	UserCommands := commands.NewUserCommands(createUserHandler, updateUserHandler, deleteUserHandler)
	UserQueries := queries.NewUserQueries(getUserByIdHandler, searchUserHandler)
	return &UserService{
		Commands: UserCommands,
		Queries:  UserQueries,
	}
}
