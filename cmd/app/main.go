package main

import (
	"github.com/tguankheng016/golang-ecommerce-monolith/config"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config/environment"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				environment.ConfigAppEnv,
				config.InitConfig,
				logger.InitLogger,
				http.NewContext,
				database.NewGormDB,
				echoServer.NewEchoServer,
			),
		),
	).Run()
}
