package main

import (
	"flag"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/server"
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
	logger.WithName("CommandService")
	s := server.NewServer(logger, cfg)
	logger.Fatal(s.Run())
}
