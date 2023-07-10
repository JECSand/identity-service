package v1

import (
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/commands"
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/JECSand/identity-service/api_gateway_service/identity/metrics"
	"github.com/JECSand/identity-service/api_gateway_service/identity/middlewares"
	"github.com/JECSand/identity-service/api_gateway_service/identity/queries"
	"github.com/JECSand/identity-service/api_gateway_service/identity/services"
	"github.com/JECSand/identity-service/pkg/constants"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/routing"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type groupsHandlers struct {
	group   *echo.Group
	log     logging.Logger
	mw      middlewares.MiddlewareManager
	cfg     *config.Config
	ps      *services.GroupService
	ms      *services.MembershipService
	v       *validator.Validate
	metrics *metrics.ApiGatewayMetrics
}

func (h *groupsHandlers) MapRoutes() {
	h.group.POST("", h.mw.RequestVerifyMiddleware(h.CreateGroup()))
	h.group.GET("/:id", h.mw.RequestVerifyMiddleware(h.GetGroupByID()))
	h.group.GET("/search", h.mw.RequestVerifyMiddleware(h.SearchGroup()))
	h.group.GET("/:id/users", h.mw.RequestVerifyMiddleware(h.GetGroupUserMemberships()))
	h.group.PUT("/:id", h.mw.RequestVerifyMiddleware(h.UpdateGroup()))
	h.group.DELETE("/:id", h.mw.RequestVerifyMiddleware(h.DeleteGroup()))
	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}

func NewGroupsHandlers(
	group *echo.Group,
	log logging.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	ps *services.GroupService,
	ms *services.MembershipService,
	v *validator.Validate,
	metrics *metrics.ApiGatewayMetrics,
) *groupsHandlers {
	return &groupsHandlers{
		group:   group,
		log:     log,
		mw:      mw,
		cfg:     cfg,
		ps:      ps,
		ms:      ms,
		v:       v,
		metrics: metrics,
	}
}

// CreateGroup
// @Tags Groups
// @Summary Create group
// @Description Create new group item
// @Accept json
// @Produce json
// @Success 201 {object} dto.CreateGroupResponseDTO
// @Router /groups [post]
func (h *groupsHandlers) CreateGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		h.metrics.CreateGroupHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "groupsHandlers.CreateGroup")
		defer span.Finish()
		createDto := &dto.CreateGroupDTO{}
		if err = c.Bind(createDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		createDto.ID, err = utilities.NewID()
		if err = h.v.StructCtx(ctx, createDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if err = h.ps.Commands.CreateGroup.Handle(ctx, commands.NewCreateGroupCommand(createDto)); err != nil {
			h.log.WarnMsg("CreateGroup", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusCreated, dto.CreateGroupResponseDTO{ID: createDto.ID})
	}
}

// GetGroupByID
// @Tags Groups
// @Summary Get group
// @Description Get group by id
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} dto.GroupResponse
// @Router /groups/{id} [get]
func (h *groupsHandlers) GetGroupByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.GetGroupByIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "groupsHandlers.GetGroupByID")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		query := queries.NewGetGroupByIdQuery(id)
		response, err := h.ps.Queries.GetGroupById.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetGroupById", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// SearchGroup
// @Tags Groups
// @Summary Search group
// @Description Get group by name with pagination
// @Accept json
// @Produce json
// @Param search query string false "search text"
// @Param page query string false "page number"
// @Param size query string false "number of elements"
// @Success 200 {object} dto.GroupsListResponse
// @Router /groups/search [get]
func (h *groupsHandlers) SearchGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.SearchGroupHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "groupsHandlers.SearchGroup")
		defer span.Finish()
		pq := utilities.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))
		query := queries.NewSearchGroupQuery(c.QueryParam(constants.Search), pq)
		response, err := h.ps.Queries.SearchGroup.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("SearchGroup", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// GetGroupUserMemberships
// @Tags Groups
// @Summary Get group user memberships
// @Description Get group user memberships by id with pagination
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param page query string false "page number"
// @Param size query string false "number of elements"
// @Success 200 {object} dto.GroupMembershipsListResponse
// @Router /groups/{id}/users [get]
func (h *groupsHandlers) GetGroupUserMemberships() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.GetUserMembershipByGroupIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "groupsHandlers.GetGroupUserMemberships")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		pq := utilities.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))
		query := queries.NewGetUserMembershipByGroupIdQuery(id, pq)
		response, err := h.ms.Queries.GetUserMembershipByGroupId.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetGroupUserMemberships", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateGroup
// @Tags Groups
// @Summary Update group
// @Description Update existing group
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} dto.UpdateGroupDTO
// @Router /groups/{id} [put]
func (h *groupsHandlers) UpdateGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.UpdateGroupHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "groupsHandlers.UpdateGroup")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		updateDto := &dto.UpdateGroupDTO{ID: id}
		if err = c.Bind(updateDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if err = h.v.StructCtx(ctx, updateDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if err = h.ps.Commands.UpdateGroup.Handle(ctx, commands.NewUpdateGroupCommand(updateDto)); err != nil {
			h.log.WarnMsg("UpdateGroup", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, updateDto)
	}
}

// DeleteGroup
// @Tags Groups
// @Summary Delete group
// @Description Delete existing group
// @Accept json
// @Produce json
// @Success 200 ""
// @Param id path string true "Group ID"
// @Router /groups/{id} [delete]
func (h *groupsHandlers) DeleteGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.DeleteGroupHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "groupsHandlers.DeleteGroup")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if err = h.ps.Commands.DeleteGroup.Handle(ctx, commands.NewDeleteGroupCommand(id)); err != nil {
			h.log.WarnMsg("DeleteGroup", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.NoContent(http.StatusOK)
	}
}

func (h *groupsHandlers) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}
