package main

import (
	"github.com/tguankheng016/golang-ecommerce-monolith/config"
	identityConfiguration "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/configurations"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	identityData "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
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
				identityData.NewUserManager,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(redis.RegisterRedisServer),
			fx.Invoke(middleware.ConfigMiddlewares),
			fx.Invoke(swagger.ConfigSwagger),
			fx.Invoke(func(gorm *gorm.DB, userManager data.IUserManager) error {
				err := database.Migrate(gorm,
					&identityModel.Role{},
					&identityModel.User{},
					&identityModel.UserToken{},
					&identityModel.UserRolePermission{},
				)
				if err != nil {
					return err
				}
				return identityData.DataSeeder(gorm, userManager)
			}),
			fx.Invoke(identityConfiguration.ConfigEndpoints),
		),
	).Run()
}
