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

type membershipsHandlers struct {
	group   *echo.Group
	log     logging.Logger
	mw      middlewares.MiddlewareManager
	cfg     *config.Config
	ps      *services.MembershipService
	v       *validator.Validate
	metrics *metrics.ApiGatewayMetrics
}

func (h *membershipsHandlers) MapRoutes() {
	h.group.POST("", h.mw.RequestVerifyMiddleware(h.CreateMembership()))
	h.group.GET("/:id", h.mw.RequestVerifyMiddleware(h.GetMembershipByID()))
	h.group.PUT("/:id", h.mw.RequestVerifyMiddleware(h.UpdateMembership()))
	h.group.DELETE("/:id", h.mw.RequestVerifyMiddleware(h.DeleteMembership()))
	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}

func NewMembershipsHandlers(
	group *echo.Group,
	log logging.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	ps *services.MembershipService,
	v *validator.Validate,
	metrics *metrics.ApiGatewayMetrics,
) *membershipsHandlers {
	return &membershipsHandlers{
		group:   group,
		log:     log,
		mw:      mw,
		cfg:     cfg,
		ps:      ps,
		v:       v,
		metrics: metrics,
	}
}

// CreateMembership
// @Tags Memberships
// @Summary Create membership
// @Description Create new membership item
// @Accept json
// @Produce json
// @Success 201 {object} dto.CreateMembershipResponseDTO
// @Router /memberships [post]
func (h *membershipsHandlers) CreateMembership() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		h.metrics.CreateMembershipHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "membershipsHandlers.CreateMembership")
		defer span.Finish()
		createDto := &dto.CreateMembershipDTO{}
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
		if err = h.ps.Commands.CreateMembership.Handle(ctx, commands.NewCreateMembershipCommand(createDto)); err != nil {
			h.log.WarnMsg("CreateMembership", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusCreated, dto.CreateMembershipResponseDTO{ID: createDto.ID})
	}
}

// GetMembershipByID
// @Tags Memberships
// @Summary Get membership
// @Description Get membership by id
// @Accept json
// @Produce json
// @Param id path string true "Membership ID"
// @Success 200 {object} dto.MembershipResponse
// @Router /memberships/{id} [get]
func (h *membershipsHandlers) GetMembershipByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.GetMembershipByIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "membershipsHandlers.GetMembershipByID")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		query := queries.NewGetMembershipByIdQuery(id)
		response, err := h.ps.Queries.GetMembershipById.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetMembershipById", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateMembership
// @Tags Memberships
// @Summary Update membership
// @Description Update existing membership
// @Accept json
// @Produce json
// @Param id path string true "Membership ID"
// @Success 200 {object} dto.UpdateMembershipDTO
// @Router /memberships/{id} [put]
func (h *membershipsHandlers) UpdateMembership() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.UpdateMembershipHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "membershipsHandlers.UpdateMembership")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		updateDto := &dto.UpdateMembershipDTO{ID: id}
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
		if err = h.ps.Commands.UpdateMembership.Handle(ctx, commands.NewUpdateMembershipCommand(updateDto)); err != nil {
			h.log.WarnMsg("UpdateMembership", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, updateDto)
	}
}

// DeleteMembership
// @Tags Memberships
// @Summary Delete membership
// @Description Delete existing membership
// @Accept json
// @Produce json
// @Success 200 ""
// @Param id path string true "Membership ID"
// @Router /memberships/{id} [delete]
func (h *membershipsHandlers) DeleteMembership() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.DeleteMembershipHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "membershipsHandlers.DeleteMembership")
		defer span.Finish()
		id, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if err = h.ps.Commands.DeleteMembership.Handle(ctx, commands.NewDeleteMembershipCommand(id)); err != nil {
			h.log.WarnMsg("DeleteMembership", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.NoContent(http.StatusOK)
	}
}

func (h *membershipsHandlers) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}
