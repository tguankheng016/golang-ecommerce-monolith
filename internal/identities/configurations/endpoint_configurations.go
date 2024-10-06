package configurations

import (
	"context"

	"github.com/labstack/echo/v4"
	authenticate "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/authenticating/v1/endpoints"
	refreshToken "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/refreshing_token/v1/endpoints"
	getting_users "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/getting_users/v1/endpoints"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"gorm.io/gorm"
)

func ConfigEndpoints(
	db *gorm.DB,
	jwtTokenGenerator jwt.IJwtTokenGenerator,
	jwtTokenValidator jwt.IJwtTokenValidator,
	log logger.ILogger,
	echo *echo.Echo,
	ctx context.Context,
) {
	getting_users.MapRoute(db, jwtTokenValidator, log, echo, ctx)
	authenticate.MapRoute(db, jwtTokenGenerator, log, echo, ctx)
	refreshToken.MapRoute(db, jwtTokenGenerator, jwtTokenValidator, log, echo, ctx)
}
