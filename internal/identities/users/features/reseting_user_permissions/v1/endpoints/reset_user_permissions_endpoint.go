package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
	"gorm.io/gorm"
)

func MapRoute(echo *echo.Echo, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/user/:userId/reset-permissions")
	group.PUT("", resetUserPermissions(permissionManager), middlewares.Authorize(permissions.PagesAdministrationUsersChangePermissions))
}

// ResetUserPermissions
// @Tags Users
// @Summary Reset user permissions
// @Description Reset user permissions
// @Accept json
// @Produce json
// @Param userId path int true "User Id"
// @Success 200
// @Security ApiKeyAuth
// @Router /api/v1/user/{userId}/reset-permissions [put]
func resetUserPermissions(permissionManager permissions.IPermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var userId int64
		if err := echo.PathParamsBinder(c).Int64("userId", &userId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := database.GetTxFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if err := tx.Model(&models.UserRolePermission{}).Where("user_id = ?", userId).Delete(&models.UserRolePermission{}).Error; err != nil && err != gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		// Commit because permission manager tx is different
		tx.Commit()

		// Reset User Permissions
		permissionManager.SetUserPermissions(ctx, userId)

		return c.NoContent(http.StatusOK)
	}
}
