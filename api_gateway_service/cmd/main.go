package main

import (
	"flag"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/controllers/access"
	"github.com/JECSand/identity-service/api_gateway_service/server"
	"github.com/JECSand/identity-service/pkg/authentication"
	"github.com/JECSand/identity-service/pkg/logging"
	"log"
)

func main() {
	flag.Parse()
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.NewAppLogger(cfg.Logger)
	logger.InitLogger()
	logger.WithName("GatewayService")
	authCfg := authentication.NewAuthConfig(1, 4380, cfg.ServiceSettings.JWTSalt)
	auth := authentication.NewAuthenticator(logger, access.DefaultAccessRules(), authCfg)
	s := server.NewServer(logger, auth, cfg)
	logger.Fatal(s.Run())
}
