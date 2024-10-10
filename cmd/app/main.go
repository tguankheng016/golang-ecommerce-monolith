package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/tguankheng016/golang-ecommerce-monolith/config"
	identityConfiguration "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/configurations"
	identityData "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config/environment"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/middleware"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/redis"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/server"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/swagger"
	"gorm.io/gorm"

	"go.uber.org/fx"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				environment.ConfigAppEnv,
				config.InitConfig,
				logger.InitLogger,
				http.NewContext,
				database.NewGormDB,
				redis.NewRedisClient,
				echoServer.NewEchoServer,
				jwt.NewJwtTokenGenerator,
				jwt.NewJwtTokenValidator,
				permissions.NewPermissionChecker,
				validator.New,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(redis.RegisterRedisServer),
			fx.Invoke(middleware.ConfigMiddlewares),
			fx.Invoke(swagger.ConfigSwagger),
			fx.Invoke(func(db *gorm.DB) error {
				if err := database.RunGooseMigration(db); err != nil {
					return err
				}
				return identityData.DataSeeder(db)
			}),
			fx.Invoke(identityConfiguration.ConfigEndpoints),
		),
	).Run()
}
