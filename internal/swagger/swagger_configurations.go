package swagger

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/tguankheng016/golang-ecommerce-monolith/docs"
)

func ConfigSwagger(e *echo.Echo) {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Golang Ecommerce Api"
	docs.SwaggerInfo.Description = "Golang Ecommerce Monolith Api"
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
