package permissions

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpServer "github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"go.uber.org/zap"
)

// SetupAuthorization returns a middleware that retrieves the granted permissions for the current user and store them in the request context.
// The middleware will first try to get the current user id from the request context. If the user id is not found, it will skip the permission check.
// The middleware will then get the granted permissions from the permission manager and store them in the request context.
// If the permission manager returns an error, the middleware will log the error and continue to the next middleware.
func SetupAuthorization(api huma.API, permissionManager IPermissionManager) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		context := ctx.Context()

		userId, ok := httpServer.GetCurrentUser(context)
		if !ok {
			next(ctx)
			return
		}

		grantedPermissions, err := permissionManager.GetGrantedPermissions(context, userId)
		if err != nil {
			logging.Logger.Error("Error when getting granted permissions: ", zap.Error(err))
		}

		ctx = httpServer.SetCurrentUserPermissions(ctx, grantedPermissions)

		next(ctx)
	}
}

func Authorize(api huma.API, permission string) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		context := ctx.Context()
		_, ok := httpServer.GetCurrentUser(context)
		if !ok {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "The current user did not log in to the application")
			return
		}

		if permission != "" {
			accessDeniedError := "You do not have permission to access this resource"

			grantedPermissions, ok := httpServer.GetCurrentUserPermissions(context)
			if !ok {
				huma.WriteErr(api, ctx, http.StatusForbidden, accessDeniedError)
				return
			}

			if _, ok := grantedPermissions[permission]; !ok {
				huma.WriteErr(api, ctx, http.StatusForbidden, accessDeniedError)
				return
			}
		}

		next(ctx)
	}
}
