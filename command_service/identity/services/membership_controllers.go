package services

import (
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/commands"
	"github.com/JECSand/identity-service/command_service/identity/queries"
	"github.com/JECSand/identity-service/command_service/identity/repositories"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
)

// MembershipService ...
type MembershipService struct {
	Commands *commands.MembershipCommands
	Queries  *queries.MembershipQueries
}

// NewMembershipService ...
func NewMembershipService(log logging.Logger, cfg *config.Config, pgRepo repositories.Repository, kafkaProducer kafkaClient.Producer) *MembershipService {
	updateMembershipHandler := commands.NewUpdateMembershipHandler(log, cfg, pgRepo, kafkaProducer)
	createMembershipHandler := commands.NewCreateMembershipHandler(log, cfg, pgRepo, kafkaProducer)
	deleteMembershipHandler := commands.NewDeleteMembershipHandler(log, cfg, pgRepo, kafkaProducer)
	getMembershipByIdHandler := queries.NewGetMembershipByIdHandler(log, cfg, pgRepo)
	getUserMembershipByIdHandler := queries.NewGetUserMembershipByIdHandler(log, cfg, pgRepo)
	getGroupMembershipByIdHandler := queries.NewGetGroupMembershipByIdHandler(log, cfg, pgRepo)
	countMembershipsHandler := queries.NewCountMembershipsHandler(log, cfg, pgRepo)
	membershipCommands := commands.NewMembershipCommands(createMembershipHandler, updateMembershipHandler, deleteMembershipHandler)
	membershipQueries := queries.NewMembershipQueries(getMembershipByIdHandler, getUserMembershipByIdHandler, getGroupMembershipByIdHandler, countMembershipsHandler)
	return &MembershipService{
		Commands: membershipCommands,
		Queries:  membershipQueries,
	}
}
