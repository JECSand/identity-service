package services

import (
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	"github.com/JECSand/identity-service/command_service/identity/queries"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
)

// GroupService ...
type GroupService struct {
	Commands *commands.GroupCommands
	Queries  *queries.GroupQueries
}

// NewGroupService ...
func NewGroupService(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *GroupService {
	updateGroupHandler := commands.NewUpdateGroupHandler(log, cfg, pgRepo, kafkaProducer)
	createGroupHandler := commands.NewCreateGroupHandler(log, cfg, pgRepo, kafkaProducer)
	deleteGroupHandler := commands.NewDeleteGroupHandler(log, cfg, pgRepo, kafkaProducer)
	getGroupByIdHandler := queries.NewGetGroupByIdHandler(log, cfg, pgRepo)
	countGroupsHandler := queries.NewCountGroupsHandler(log, cfg, pgRepo)
	GroupCommands := commands.NewGroupCommands(createGroupHandler, updateGroupHandler, deleteGroupHandler)
	GroupQueries := queries.NewGroupQueries(getGroupByIdHandler, countGroupsHandler)
	return &GroupService{
		Commands: GroupCommands,
		Queries:  GroupQueries,
	}
}
