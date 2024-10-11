package configurations

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	authenticate "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/authenticating/v1/endpoints"
	refreshToken "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/refreshing_token/v1/endpoints"
	sign_out "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/features/signing_out/v1/endpoints"
	creating_role "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/roles/features/creating_role/v1/endpoints"
	deleting_role "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/roles/features/deleting_role/v1/endpoints"
	get_role_by_id "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/roles/features/get_role_by_id/v1/endpoints"
	getting_roles "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/roles/features/getting_roles/v1/endpoints"
	updating_role "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/roles/features/updating_role/v1/endpoints"
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
	permissionManager permissions.IPermissionManager,
	validator *validator.Validate,
	log logger.ILogger,
	echo *echo.Echo,
) {
	// Users
	getting_users.MapRoute(echo, validator, jwtTokenValidator, permissionManager)
	get_user_by_id.MapRoute(echo, jwtTokenValidator, permissionManager)
	creating_user.MapRoute(echo, validator, jwtTokenValidator, permissionManager)
	updating_user.MapRoute(echo, validator, jwtTokenValidator, permissionManager)
	deleting_user.MapRoute(echo, jwtTokenValidator, permissionManager)

	// Roles
	getting_roles.MapRoute(echo, validator, jwtTokenValidator, permissionManager)
	get_role_by_id.MapRoute(echo, jwtTokenValidator, permissionManager)
	creating_role.MapRoute(echo, validator, jwtTokenValidator, permissionManager)
	updating_role.MapRoute(echo, validator, jwtTokenValidator, permissionManager)
	deleting_role.MapRoute(echo, jwtTokenValidator, permissionManager)

	// Accounts
	authenticate.MapRoute(echo, validator, jwtTokenGenerator)
	refreshToken.MapRoute(echo, validator, jwtTokenGenerator, jwtTokenValidator)
	sign_out.MapRoute(echo, jwtTokenValidator, jwtTokenGenerator)
}
