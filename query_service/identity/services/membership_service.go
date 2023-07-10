package services

import (
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/cache"
	"github.com/JECSand/identity-service/query_service/identity/data"
	"github.com/JECSand/identity-service/query_service/identity/events"
	"github.com/JECSand/identity-service/query_service/identity/queries"
)

type MembershipService struct {
	Events  *events.MembershipEvents
	Queries *queries.MembershipQueries
}

func NewMembershipService(
	log logging.Logger,
	cfg *config.Config,
	mongoDB data.Database,
	redisCache cache.Cache,
) *MembershipService {
	createMembershipHandler := events.NewCreateMembershipEventHandler(log, cfg, mongoDB, redisCache)
	deleteMembershipEventHandler := events.NewDeleteMembershipEventHandler(log, cfg, mongoDB, redisCache)
	updateMembershipEventHandler := events.NewUpdateMembershipEventHandler(log, cfg, mongoDB, redisCache)
	getMembershipByIdHandler := queries.NewGetMembershipByIdHandler(log, cfg, mongoDB, redisCache)
	getGroupMembershipHandler := queries.NewGetGroupMembershipHandler(log, cfg, mongoDB, redisCache)
	getUserMembershipHandler := queries.NewGetUserMembershipHandler(log, cfg, mongoDB, redisCache)
	membershipEvents := events.NewMembershipEvents(createMembershipHandler, updateMembershipEventHandler, deleteMembershipEventHandler)
	membershipQueries := queries.NewMembershipQueries(getMembershipByIdHandler, getGroupMembershipHandler, getUserMembershipHandler)
	return &MembershipService{
		Events:  membershipEvents,
		Queries: membershipQueries,
	}
}
