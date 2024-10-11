package middleware

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
	"gorm.io/gorm"
)

func ConfigMiddlewares(
	e *echo.Echo,
	db *gorm.DB,
	validator *validator.Validate,
	logger logger.ILogger,
	jwtTokenValidator jwt.IJwtTokenValidator,
	permissionManager permissions.IPermissionManager,
) {
	skipper := func(c echo.Context) bool {
		return strings.Contains(c.Request().URL.Path, "swagger") ||
			strings.Contains(c.Request().URL.Path, "metrics") ||
			strings.Contains(c.Request().URL.Path, "health") ||
			strings.Contains(c.Request().URL.Path, "favicon.ico")
	}

	e.HideBanner = false

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:   5,
		Skipper: skipper,
	}))

	e.Use(middleware.BodyLimit("2M"))

	e.Use(middlewares.SetupAuthenticate(skipper, jwtTokenValidator, logger))
	e.Use(middlewares.SetupAuthorize(skipper, permissionManager, logger))
	e.Use(middlewares.SetupTransaction(skipper, db))
}
