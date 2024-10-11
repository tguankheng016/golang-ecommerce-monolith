package middlewares

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func Authorize(permissionManager permissions.IPermissionManager, permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if permission == "" {
				return next(c)
			}

			// Ignore check permission in test
			env := os.Getenv("APP_ENV")
			if env == "test" {
				return next(c)
			}

			userId, ok := echoServer.GetCurrentUser(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("get user id error"))
			}

			ctx := c.Request().Context()

			isGranted, err := permissionManager.IsGranted(ctx, userId, permission)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			if !isGranted {
				return echo.NewHTTPError(http.StatusForbidden, errors.New("permission denied"))
			}

			return next(c)
		}
	}
}
