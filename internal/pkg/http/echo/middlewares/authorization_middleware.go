package middlewares

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func Authorize(checker permissions.IPermissionChecker, permission string) echo.MiddlewareFunc {
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

			userId, ok := c.Get("userId").(int64)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("get user id error"))
			}

			isGranted, err := checker.IsGranted(userId, permission)
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
