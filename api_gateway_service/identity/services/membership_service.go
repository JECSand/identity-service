package services

import (
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/commands"
	"github.com/JECSand/identity-service/api_gateway_service/identity/queries"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	membershipQueryService "github.com/JECSand/identity-service/query_service/protos/membership_query"
)

type MembershipService struct {
	Commands *commands.MembershipCommands
	Queries  *queries.MembershipQueries
}

func NewMembershipService(log logging.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer, rsClient membershipQueryService.MembershipQueryServiceClient) *MembershipService {
	createMembershipHandler := commands.NewCreateMembershipHandler(log, cfg, kafkaProducer)
	updateMembershipHandler := commands.NewUpdateMembershipHandler(log, cfg, kafkaProducer)
	deleteMembershipHandler := commands.NewDeleteMembershipHandler(log, cfg, kafkaProducer)
	getMembershipByIdHandler := queries.NewGetMembershipByIdHandler(log, cfg, rsClient)
	getUserMembershipByGroupIdHandler := queries.NewGetUserMembershipByGroupIHandler(log, cfg, rsClient)
	getGroupMembershipByUserIdHandler := queries.NewGetGroupMembershipByUserIdHandler(log, cfg, rsClient)
	MembershipCommands := commands.NewMembershipCommands(createMembershipHandler, updateMembershipHandler, deleteMembershipHandler)
	MembershipQueries := queries.NewMembershipQueries(getMembershipByIdHandler, getUserMembershipByGroupIdHandler, getGroupMembershipByUserIdHandler)
	return &MembershipService{
		Commands: MembershipCommands,
		Queries:  MembershipQueries,
	}
}
