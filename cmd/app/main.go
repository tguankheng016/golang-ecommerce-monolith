package main

import (
	"github.com/tguankheng016/golang-ecommerce-monolith/config"
	identityConfiguration "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/configurations"
	identitySeed "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	identityMapping "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/mappings"
	identityModel "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/middleware"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config/environment"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/redis"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/server"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/swagger"
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
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(redis.RegisterRedisServer),
			fx.Invoke(middleware.ConfigMiddlewares),
			fx.Invoke(swagger.ConfigSwagger),
			fx.Invoke(func(gorm *gorm.DB) error {
				err := database.Migrate(gorm,
					&identityModel.Role{},
					&identityModel.User{},
					&identityModel.UserToken{},
					&identityModel.UserRolePermission{},
				)
				if err != nil {
					return err
				}
				return identitySeed.DataSeeder(gorm)
			}),
			fx.Invoke(identityMapping.ConfigureMappings),
			fx.Invoke(identityConfiguration.ConfigEndpoints),
		),
	).Run()
}
