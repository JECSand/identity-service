package services

import (
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	"github.com/JECSand/identity-service/command_service/identity/queries"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
)

// AuthService ...
type AuthService struct {
	Commands *commands.AuthCommands
	Queries  *queries.AuthQueries
}

// NewAuthService ...
func NewAuthService(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *AuthService {
	blacklistTokenHandler := commands.NewBlacklistTokenHandler(log, cfg, pgRepo, kafkaProducer)
	passwordUpdateHandler := commands.NewUpdatePasswordHandler(log, cfg, pgRepo, kafkaProducer)
	checkBlacklistHandler := queries.NewCheckTokenBlacklistHandler(log, cfg, pgRepo)
	userCommands := commands.NewAuthCommands(blacklistTokenHandler, passwordUpdateHandler)
	userQueries := queries.NewAuthQueries(checkBlacklistHandler)
	return &AuthService{
		Commands: userCommands,
		Queries:  userQueries,
	}
}
