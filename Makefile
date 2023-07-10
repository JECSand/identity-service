.PHONY:

run_gateway_service:
	go run api_gateway_service/cmd/main.go -config=./api_gateway_service/config/config.yaml

run_command_service:
	go run command_service/cmd/main.go -config=./command_service/config/config.yaml

run_query_service:
	go run query_service/cmd/main.go -config=./query_service/config/config.yaml

# ==============================================================================
# Docker

docker_dev:
	@echo Starting local docker dev compose
	docker-compose -f docker-compose.yaml up --build

local:
	@echo Starting local docker compose
	docker-compose -f docker-compose.local.yaml up -d --build


# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)

# ==============================================================================
# Modules support

tidy:
	go mod tidy

deps-reset:
	git checkout -- go.mod
	go mod tidy

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

deps-cleancache:
	go clean -modcache


# ==============================================================================
# Make local SSL Certificate

cert:
	echo "Generating SSL certificates"
	cd ./ssl && sh generate.sh

# ==============================================================================
# Go migrate postgresql https://github.com/golang-migrate/migrate

DB_NAME = users
DB_HOST = localhost
DB_PORT = 5432
SSL_MODE = disable

force_db:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path migrations force 1

version_db:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path migrations down 1


# ==============================================================================
# MongoDB

mongo:
	cd ./migrations && mongo admin -u admin -p admin < initDB.js


# ==============================================================================
# Swagger

swagger:
	@echo Starting swagger generating
	swag init -g **/**/*.go



# ==============================================================================
# Proto

proto_kafka:
	@echo Generating kafka proto
	cd protos/kafka && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. kafka.proto

proto_command:
	@echo Generating user command service proto
	cd command_service/proto/user_command && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. user_command.proto

proto_command_message:
	@echo Generating user command messages service proto
	cd command_service/proto/user_command && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. user_command_messages.proto

proto_query:
	@echo Generating user query service proto
	cd query_service/protos/user_query && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. user_query.proto

proto_query_message:
	@echo Generating user query messages service proto
	cd query_service/protos/user_query && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. user_query_messages.proto