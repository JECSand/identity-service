package main

import (
	"flag"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/server"
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
	logger.WithName("QueryService")
	s := server.NewServer(logger, cfg)
	logger.Fatal(s.Run())
}
