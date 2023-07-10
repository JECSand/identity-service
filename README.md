# identity-service
An open sourced identity service written in Go and leveraging clean architecture + CQRS principles.

## Getting Started
### Local
1. Run docker-compose.local.yaml.
```shell
$ make local
```
2. Use migration files to setup MongoDB and Postgres databases.
3. Manually start the services via the makefile:
```shell
$ make run_query_service
$ make run_command_service
$ make run_gateway_service
```

### Development
1. Run docker-compose.yaml.
```shell
$ make docker_dev
```