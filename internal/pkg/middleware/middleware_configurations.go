package middleware

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"gorm.io/gorm"
)

func ConfigMiddlewares(e *echo.Echo, db *gorm.DB, validator *validator.Validate) {
	e.HideBanner = false

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Use(middleware.BodyLimit("2M"))
	e.Use(middlewares.TransactionalContextMiddleware(db))
}
