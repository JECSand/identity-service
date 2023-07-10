package services

import (
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/commands"
	"github.com/JECSand/identity-service/api_gateway_service/identity/queries"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	authQueryService "github.com/JECSand/identity-service/query_service/protos/auth_query"
)

type AuthService struct {
	Commands *commands.AuthCommands
	Queries  *queries.AuthQueries
}

func NewAuthService(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer, rsClient authQueryService.AuthQueryServiceClient) *AuthService {
	blacklistTokenHandler := commands.NewBlacklistTokenHandler(log, cfg, kafkaProducer)
	passwordUpdateHandler := commands.NewUpdatePasswordHandler(log, cfg, kafkaProducer)
	authenticateHandler := queries.NewAuthenticateHandler(log, cfg, rsClient)
	validateHandler := queries.NewValidateHandler(log, cfg, rsClient)
	AuthCommands := commands.NewAuthCommands(blacklistTokenHandler, passwordUpdateHandler)
	AuthQueries := queries.NewAuthQueries(authenticateHandler, validateHandler)
	return &AuthService{
		Commands: AuthCommands,
		Queries:  AuthQueries,
	}
}
