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

type usersHandlers struct {
	group   *echo.Group
	log     logging.Logger
	mw      middlewares.MiddlewareManager
	cfg     *config.Config
	ps      *services.UserService
	ms      *services.MembershipService
	v       *validator.Validate
	metrics *metrics.ApiGatewayMetrics
}

func (h *usersHandlers) MapRoutes() {
	h.group.POST("", h.mw.RequestVerifyMiddleware(h.CreateUser()))
	h.group.GET("/:id", h.mw.RequestVerifyMiddleware(h.GetUserByID()))
	h.group.GET("/search", h.mw.RequestVerifyMiddleware(h.SearchUser()))
	h.group.GET("/:id/groups", h.mw.RequestVerifyMiddleware(h.GetUserGroupMemberships()))
	h.group.PUT("/:id", h.mw.RequestVerifyMiddleware(h.UpdateUser()))
	h.group.DELETE("/:id", h.mw.RequestVerifyMiddleware(h.DeleteUser()))
	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}

func NewUsersHandlers(
	group *echo.Group,
	log logging.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	ps *services.UserService,
	ms *services.MembershipService,
	v *validator.Validate,
	metrics *metrics.ApiGatewayMetrics,
) *usersHandlers {
	return &usersHandlers{
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

// CreateUser
// @Tags Users
// @Summary Create user
// @Description Create new user item
// @Accept json
// @Produce json
// @Success 201 {object} dto.CreateUserResponseDTO
// @Router /users [post]
func (h *usersHandlers) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		h.metrics.CreateUserHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "usersHandlers.CreateUser")
		defer span.Finish()
		createDto := &dto.CreateUserDTO{}
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
		if err = h.ps.Commands.CreateUser.Handle(ctx, commands.NewCreateUserCommand(createDto)); err != nil {
			h.log.WarnMsg("CreateUser", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusCreated, dto.CreateUserResponseDTO{ID: createDto.ID})
	}
}

// GetUserByID
// @Tags Users
// @Summary Get user
// @Description Get user by id
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Router /users/{id} [get]
func (h *usersHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.GetUserByIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "usersHandlers.GetUserByID")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		query := queries.NewGetUserByIdQuery(id)
		response, err := h.ps.Queries.GetUserById.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetUserById", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// SearchUser
// @Tags Users
// @Summary Search user
// @Description Get user by name with pagination
// @Accept json
// @Produce json
// @Param search query string false "search text"
// @Param page query string false "page number"
// @Param size query string false "number of elements"
// @Success 200 {object} dto.UsersListResponse
// @Router /users/search [get]
func (h *usersHandlers) SearchUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.SearchUserHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "usersHandlers.SearchUser")
		defer span.Finish()
		pq := utilities.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))
		query := queries.NewSearchUserQuery(c.QueryParam(constants.Search), pq)
		response, err := h.ps.Queries.SearchUser.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("SearchUser", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// GetUserGroupMemberships
// @Tags Users
// @Summary Get user group memberships
// @Description Get user group memberships by id with pagination
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param page query string false "page number"
// @Param size query string false "number of elements"
// @Success 200 {object} dto.GroupMembershipsListResponse
// @Router /users/{id}/groups [get]
func (h *usersHandlers) GetUserGroupMemberships() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.GetGroupMembershipByUserIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "usersHandlers.GetUserGroupMemberships")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		pq := utilities.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))
		query := queries.NewGetGroupMembershipByUserIdQuery(id, pq)
		response, err := h.ms.Queries.GetGroupMembershipByUserId.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetUserGroupMemberships", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateUser
// @Tags Users
// @Summary Update user
// @Description Update existing user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UpdateUserDTO
// @Router /users/{id} [put]
func (h *usersHandlers) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.UpdateUserHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "usersHandlers.UpdateUser")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		updateDto := &dto.UpdateUserDTO{ID: id}
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
		if err = h.ps.Commands.UpdateUser.Handle(ctx, commands.NewUpdateUserCommand(updateDto)); err != nil {
			h.log.WarnMsg("UpdateUser", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, updateDto)
	}
}

// DeleteUser
// @Tags Users
// @Summary Delete user
// @Description Delete existing user
// @Accept json
// @Produce json
// @Success 200 ""
// @Param id path string true "User ID"
// @Router /users/{id} [delete]
func (h *usersHandlers) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.DeleteUserHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "usersHandlers.DeleteUser")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if err = h.ps.Commands.DeleteUser.Handle(ctx, commands.NewDeleteUserCommand(id)); err != nil {
			h.log.WarnMsg("DeleteUser", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.NoContent(http.StatusOK)
	}
}

func (h *usersHandlers) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}
