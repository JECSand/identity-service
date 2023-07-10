package middlewares

import (
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/JECSand/identity-service/api_gateway_service/identity/dto"
	"github.com/JECSand/identity-service/api_gateway_service/identity/queries"
	"github.com/JECSand/identity-service/api_gateway_service/identity/services"
	"github.com/JECSand/identity-service/pkg/authentication"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	RequestVerifyMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log  logging.Logger
	auth authentication.Authenticator
	cfg  *config.Config
	as   *services.AuthService
}

func NewMiddlewareManager(log logging.Logger, auth authentication.Authenticator, cfg *config.Config, as *services.AuthService) *middlewareManager {
	return &middlewareManager{
		log:  log,
		auth: auth,
		cfg:  cfg,
		as:   as,
	}
}

// RequestVerifyMiddleware ...
func (mw *middlewareManager) RequestVerifyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		session, err := mw.auth.AuthorizeREST(req, req.Method+" "+req.URL.String())
		if err != nil {
			mw.log.WarnMsg("auth.AuthorizeREST", err)
			return ctx.JSON(http.StatusUnauthorized, dto.ErrorDTO{Message: err.Error()})
		}
		query := queries.NewValidateQuery(session.UserId, req.Header.Get("Authorization"), enums.TOKEN)
		val, err := mw.as.Queries.Validate.Handle(req.Context(), query)
		if err != nil {
			mw.log.WarnMsg("as.Queries.Validate.Handle", err)
			return ctx.JSON(http.StatusUnauthorized, dto.ErrorDTO{Message: err.Error()})
		}
		if val.Status != 200 {
			return ctx.JSON(http.StatusUnauthorized, dto.ErrorDTO{Message: "unauthorized"})
		}
		return next(ctx)
	}
}

func (mw *middlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)
		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start)
		if !mw.checkIgnoredURI(ctx.Request().RequestURI, mw.cfg.Http.IgnoreLogUrls) {
			mw.log.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, s)
		}
		return err
	}
}

func (mw *middlewareManager) checkIgnoredURI(requestURI string, uriList []string) bool {
	for _, s := range uriList {
		if strings.Contains(requestURI, s) {
			return true
		}
	}
	return false
}
