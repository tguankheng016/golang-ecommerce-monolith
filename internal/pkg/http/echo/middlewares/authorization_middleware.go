package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

// SetupAuthorize returns an Echo middleware that sets the user's permissions in the context, by calling
// the SetUserPermissions method of the provided permissionManager. The granted permissions are then
// available in the context using the echoServer.GetCurrentUserPermissions function.
//
// The middleware uses the provided skipper to determine if the request should be skipped. If the skipper
// returns true, the middleware does not call the permissionManager and instead calls the next handler in
// the chain.
//
// If the permissionManager returns an error, the middleware logs the error using the provided logger, but
// does not return an error to the client.
func SetupAuthorize(skipper echoMiddleware.Skipper, permissionManager permissions.IPermissionManager, logger logger.ILogger) echo.MiddlewareFunc {
	// Defaults
	if skipper == nil {
		skipper = echoMiddleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			userId, ok := echoServer.GetCurrentUser(c)
			if !ok {
				return next(c)
			}

			ctx := c.Request().Context()

			grantedPermissions, err := permissionManager.SetUserPermissions(ctx, userId)
			if err != nil {
				logger.Error(err)
			}

			echoServer.SetCurrentUserPermissions(c, grantedPermissions)

			return next(c)
		}
	}
}

// Authorize returns an Echo middleware that enforces that the current user has the specified permission.
// If the current user does not have the specified permission, the middleware returns an HTTP 403 error.
// If no permission is specified, the middleware perform any checks.
func Authorize(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, ok := echoServer.GetCurrentUser(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("The current user did not log in to the application"))
			}

			accessDeniedError := errors.New("You do not have permission to access this resource")

			grantedPermissions, ok := echoServer.GetCurrentUserPermissions(c)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, accessDeniedError)
			}

			if permission != "" {
				if _, ok := grantedPermissions[permission]; !ok {
					return echo.NewHTTPError(http.StatusForbidden, accessDeniedError)
				}
			}

			return next(c)
		}
	}
}
