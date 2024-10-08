package configurations

import (
	"context"

	"github.com/labstack/echo/v4"
	authenticate "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/authenticating/v1/endpoints"
	refreshToken "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/refreshing_token/v1/endpoints"
	creating_user "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/creating_user/v1/endpoints"
	getting_users "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/getting_users/v1/endpoints"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func ConfigEndpoints(
	jwtTokenGenerator jwt.IJwtTokenGenerator,
	jwtTokenValidator jwt.IJwtTokenValidator,
	checker permissions.IPermissionChecker,
	log logger.ILogger,
	echo *echo.Echo,
	ctx context.Context,
) {
	getting_users.MapRoute(jwtTokenValidator, checker, log, echo, ctx)
	creating_user.MapRoute(jwtTokenValidator, checker, log, echo, ctx)
	authenticate.MapRoute(jwtTokenGenerator, log, echo, ctx)
	refreshToken.MapRoute(jwtTokenGenerator, jwtTokenValidator, log, echo, ctx)
}
