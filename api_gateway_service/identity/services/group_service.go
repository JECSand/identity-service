package services

import (
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/commands"
	"github.com/JECSand/identity-service/api_gateway_service/identity/queries"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	groupQueryService "github.com/JECSand/identity-service/query_service/protos/group_query"
)

type GroupService struct {
	Commands *commands.GroupCommands
	Queries  *queries.GroupQueries
}

func NewGroupService(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer, rsClient groupQueryService.GroupQueryServiceClient) *GroupService {
	createGroupHandler := commands.NewCreateGroupHandler(log, cfg, kafkaProducer)
	updateGroupHandler := commands.NewUpdateGroupHandler(log, cfg, kafkaProducer)
	deleteGroupHandler := commands.NewDeleteGroupHandler(log, cfg, kafkaProducer)
	getGroupByIdHandler := queries.NewGetGroupByIdHandler(log, cfg, rsClient)
	searchGroupHandler := queries.NewSearchGroupHandler(log, cfg, rsClient)
	GroupCommands := commands.NewGroupCommands(createGroupHandler, updateGroupHandler, deleteGroupHandler)
	GroupQueries := queries.NewGroupQueries(getGroupByIdHandler, searchGroupHandler)
	return &GroupService{
		Commands: GroupCommands,
		Queries:  GroupQueries,
	}
}
