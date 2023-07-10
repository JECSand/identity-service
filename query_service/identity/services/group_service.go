package services

import (
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/cache"
	"github.com/JECSand/identity-service/query_service/identity/data"
	"github.com/JECSand/identity-service/query_service/identity/events"
	"github.com/JECSand/identity-service/query_service/identity/queries"
)

type GroupService struct {
	Events  *events.GroupEvents
	Queries *queries.GroupQueries
}

func NewGroupService(
	log logging.Logger,
	cfg *config.Config,
	mongoDB data.Database,
	redisCache cache.Cache,
) *GroupService {
	createGroupHandler := events.NewCreateGroupEventHandler(log, cfg, mongoDB, redisCache)
	deleteGroupEventHandler := events.NewDeleteGroupEventHandler(log, cfg, mongoDB, redisCache)
	updateGroupEventHandler := events.NewUpdateGroupEventHandler(log, cfg, mongoDB, redisCache)
	getGroupByIdHandler := queries.NewGetGroupByIdHandler(log, cfg, mongoDB, redisCache)
	searchGroupHandler := queries.NewSearchGroupHandler(log, cfg, mongoDB, redisCache)
	groupEvents := events.NewGroupEvents(createGroupHandler, updateGroupEventHandler, deleteGroupEventHandler)
	groupQueries := queries.NewGroupQueries(getGroupByIdHandler, searchGroupHandler)
	return &GroupService{
		Events:  groupEvents,
		Queries: groupQueries,
	}
}
