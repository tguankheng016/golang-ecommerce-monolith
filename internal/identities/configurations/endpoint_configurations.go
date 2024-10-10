package configurations

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	authenticate "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/authenticating/v1/endpoints"
	refreshToken "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/refreshing_token/v1/endpoints"
	creating_user "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/creating_user/v1/endpoints"
	deleting_user "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/deleting_user/v1/endpoints"
	get_user_by_id "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/get_user_by_id/v1/endpoints"
	getting_users "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/getting_users/v1/endpoints"
	updating_user "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/updating_user/v1/endpoints"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func ConfigEndpoints(
	jwtTokenGenerator jwt.IJwtTokenGenerator,
	jwtTokenValidator jwt.IJwtTokenValidator,
	checker permissions.IPermissionChecker,
	validator *validator.Validate,
	log logger.ILogger,
	echo *echo.Echo,
) {
	// Users
	getting_users.MapRoute(echo, validator, jwtTokenValidator, checker)
	get_user_by_id.MapRoute(echo, jwtTokenValidator, checker)
	creating_user.MapRoute(echo, validator, jwtTokenValidator, checker)
	updating_user.MapRoute(echo, validator, jwtTokenValidator, checker)
	deleting_user.MapRoute(echo, jwtTokenValidator, checker)

	// Accounts
	authenticate.MapRoute(echo, validator, jwtTokenGenerator)
	refreshToken.MapRoute(echo, validator, jwtTokenGenerator, jwtTokenValidator)
}
