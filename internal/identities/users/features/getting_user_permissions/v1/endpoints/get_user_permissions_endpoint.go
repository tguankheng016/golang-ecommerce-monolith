package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/core/helpers"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

type UserPermissionsResult struct {
	Items []string `json:"items"`
} // @name UserPermissionsResult

func MapRoute(echo *echo.Echo, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/user/:userId/permissions")
	group.GET("", getUserPermissions(permissionManager), middlewares.Authorize(permissions.PagesAdministrationUsersChangePermissions))
}

// GetUserPermissions
// @Tags Users
// @Summary Get user permissions
// @Description Get user permissions
// @Accept json
// @Produce json
// @Param userId path int true "User Id"
// @Success 200 {object} UserPermissionsResult
// @Security ApiKeyAuth
// @Router /api/v1/user/{userId}/permissions [get]
func getUserPermissions(permissionManager permissions.IPermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var userId int64
		if err := echo.PathParamsBinder(c).Int64("userId", &userId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		grantedPermissions, err := permissionManager.SetUserPermissions(ctx, userId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, UserPermissionsResult{Items: helpers.MapKeysToSlice(grantedPermissions)})
	}
}
