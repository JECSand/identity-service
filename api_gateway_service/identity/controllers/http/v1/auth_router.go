package v1

import (
	"github.com/JECSand/identity-service/api_gateway_service/config"
	commands2 "github.com/JECSand/identity-service/api_gateway_service/identity/commands"
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/JECSand/identity-service/api_gateway_service/identity/metrics"
	"github.com/JECSand/identity-service/api_gateway_service/identity/middlewares"
	"github.com/JECSand/identity-service/api_gateway_service/identity/queries"
	services2 "github.com/JECSand/identity-service/api_gateway_service/identity/services"
	"github.com/JECSand/identity-service/pkg/authentication"
	"github.com/JECSand/identity-service/pkg/enums"
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

type authHandlers struct {
	group   *echo.Group
	log     logging.Logger
	auth    authentication.Authenticator
	mw      middlewares.MiddlewareManager
	cfg     *config.Config
	as      *services2.AuthService
	us      *services2.UserService
	v       *validator.Validate
	metrics *metrics.ApiGatewayMetrics
}

func (h *authHandlers) MapRoutes() {
	h.group.POST("", h.Authenticate())
	h.group.GET("", h.mw.RequestVerifyMiddleware(h.Validate()))
	h.group.DELETE("", h.mw.RequestVerifyMiddleware(h.Invalidate()))
	h.group.POST("/password", h.mw.RequestVerifyMiddleware(h.UpdatePassword()))
	h.group.POST("/register", h.Register())
	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}

func NewAuthHandlers(
	group *echo.Group,
	log logging.Logger,
	auth authentication.Authenticator,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	as *services2.AuthService,
	v *validator.Validate,
	metrics *metrics.ApiGatewayMetrics,
) *authHandlers {
	return &authHandlers{
		group:   group,
		log:     log,
		auth:    auth,
		mw:      mw,
		cfg:     cfg,
		as:      as,
		v:       v,
		metrics: metrics,
	}
}

// Authenticate
// @Tags Auth
// @Summary Authenticate
// @Description Authenticates a user based on credentials
// @Accept json
// @Produce json
// @Success 200 {object} dto.AuthenticateResponse
// @Router /auth [post]
func (h *authHandlers) Authenticate() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		h.metrics.AuthenticateHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "authHandlers.Authenticate")
		defer span.Finish()
		authDto := &dto.AuthenticateDTO{}
		if err = c.Bind(authDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if err = h.v.StructCtx(ctx, authDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		query := queries.NewAuthenticateQuery(authDto.Email, authDto.Password)
		response, err := h.as.Queries.Authenticate.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("Authenticate", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		session := h.auth.NewSession(response.User.ID, response.User.Root, enums.USER)
		token, err := session.NewToken()
		if err != nil {
			h.log.WarnMsg("session.NewToken", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		c.Response().Header().Set("Authorization", token)
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// Register
// @Tags Auth
// @Summary Register
// @Description Register a new User
// @Accept json
// @Produce json
// @Success 200 {object} dto.CreateUserResponseDTO
// @Router /auth/register [post]
func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.RegisterHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "authHandlers.Register")
		defer span.Finish()
		var err error
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
		command := commands2.NewCreateUserCommand(createDto)
		if err = h.us.Commands.CreateUser.Handle(ctx, command); err != nil {
			h.log.WarnMsg("Invalidate", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		session := h.auth.NewSession(createDto.ID.String(), false, enums.USER)
		token, err := session.NewToken()
		if err != nil {
			h.log.WarnMsg("session.NewToken", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		c.Response().Header().Set("Authorization", token)
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, dto.CreateUserResponseDTO{ID: createDto.ID})
	}
}

// UpdatePassword
// @Tags Auth
// @Summary UpdatePassword
// @Description Updates a user's password
// @Accept json
// @Produce json
// @Success 200 {object} dto.ValidateResponse
// @Router /auth/password [post]
func (h *authHandlers) UpdatePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.UpdatePasswordHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "authHandlers.UpdatePassword")
		defer span.Finish()
		req := c.Request()
		session, err := h.auth.GetTokenSession(req.Header.Get("Authorization"))
		if err != nil {
			h.log.WarnMsg("GetTokenSession", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		updateDto := &dto.UpdatePasswordDTO{}
		if err = c.Bind(updateDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		query := queries.NewValidateQuery(session.UserId, updateDto.CurrentPassword, enums.PASSWORD)
		response, err := h.as.Queries.Validate.Handle(ctx, query)
		if err != nil || response.Status != 200 {
			h.log.WarnMsg("Validate", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		updateDto.ID, err = uuid.FromString(session.UserId)
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if err = h.v.StructCtx(ctx, updateDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		command := commands2.NewUpdatePasswordCommand(updateDto)
		if err = h.as.Commands.UpdatePassword.Handle(ctx, command); err != nil {
			h.log.WarnMsg("Invalidate", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, updateDto.ID)
	}
}

// Invalidate
// @Tags Auth
// @Summary Invalidate
// @Description Invalidates a user session
// @Accept json
// @Produce json
// @Success 200 {object} dto.invalidateDto
// @Router /auth [delete]
func (h *authHandlers) Invalidate() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		h.metrics.InvalidateHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "authHandlers.Invalidate")
		defer span.Finish()
		req := c.Request()
		invalidateDto := &dto.BlacklistTokenDTO{AccessToken: req.Header.Get("Authorization")}
		invalidateDto.ID, err = utilities.NewID()
		if err = h.v.StructCtx(ctx, invalidateDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		command := commands2.NewBlacklistTokenCommand(invalidateDto)
		if err = h.as.Commands.BlacklistToken.Handle(ctx, command); err != nil {
			h.log.WarnMsg("Invalidate", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, invalidateDto)
	}
}

// Validate
// @Tags Auth
// @Summary Validate
// @Description Validates a user session
// @Accept json
// @Produce json
// @Success 200 {object} dto.ValidateResponse
// @Router /auth [get]
func (h *authHandlers) Validate() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.ValidateHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "authHandlers.Validate")
		defer span.Finish()
		req := c.Request()
		session, err := h.auth.GetTokenSession(req.Header.Get("Authorization"))
		if err != nil {
			h.log.WarnMsg("GetTokenSession", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		query := queries.NewValidateQuery(session.UserId, req.Header.Get("Authorization"), enums.TOKEN)
		response, err := h.as.Queries.Validate.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("Validate", err)
			h.metrics.ErrorHttpRequests.Inc()
			return routing.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

func (h *authHandlers) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}
