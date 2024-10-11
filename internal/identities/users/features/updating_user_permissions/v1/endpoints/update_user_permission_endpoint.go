package endpoints

import (
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
	"gorm.io/gorm"
)

type UpdateUserPermissionDto struct {
	GrantedPermissions []string `json:"grantedPermissions"`
} // @name UpdateUserPermissionDto

func MapRoute(echo *echo.Echo, jwt jwt.IJwtTokenValidator, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/user/:userId/permissions")
	group.PUT("", updateUserPermissions(permissionManager), middlewares.ValidateToken(jwt), middlewares.Authorize(permissionManager, permissions.PagesAdministrationUsersChangePermissions))
}

// UpdateUserPermissions
// @Tags Users
// @Summary Update user permissions
// @Description Update user permissions
// @Accept json
// @Produce json
// @Param userId path int true "User Id"
// @Param UpdateUserPermissionDto body UpdateUserPermissionDto false "UpdateUserPermissionDto"
// @Success 200
// @Security ApiKeyAuth
// @Router /api/v1/user/{userId}/permissions [put]
func updateUserPermissions(permissionManager permissions.IPermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var userId int64
		if err := echo.PathParamsBinder(c).Int64("userId", &userId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var updateUserPermissionDto UpdateUserPermissionDto
		if err := c.Bind(&updateUserPermissionDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		if err := permissions.ValidatePermissionName(updateUserPermissionDto.GrantedPermissions); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		oldPermissions, err := permissionManager.SetUserPermissions(ctx, userId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		tx, err := database.RetrieveTxCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		// Prohibit
		for _, oldPermission := range oldPermissions {
			if !slices.Contains(updateUserPermissionDto.GrantedPermissions, oldPermission) {
				if err := tx.Where("user_id = ? AND name = ?", userId, oldPermission).Delete(&models.UserRolePermission{}).Error; err != nil && err != gorm.ErrRecordNotFound {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}

				// Check role got or granted or not
				ok, err := permissionManager.IsGranted(ctx, userId, oldPermission)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				if !ok {
					continue
				}

				// Prohibit at user level if role is granted
				if err := tx.Create(&models.UserRolePermission{UserId: userId, Name: oldPermission, IsGranted: false}).Error; err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
			}
		}

		// Granted
		for _, newPermission := range updateUserPermissionDto.GrantedPermissions {
			if _, ok := oldPermissions[newPermission]; !ok {
				// Check and delete any false granted user level permission
				if err := tx.Where("user_id = ? AND name = ? AND is_granted = ?", userId, newPermission, false).Delete(&models.UserRolePermission{}).Error; err != nil && err != gorm.ErrRecordNotFound {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}

				// Check role got or granted or not
				// Skip if role already have that permission
				ok, err := permissionManager.IsGranted(ctx, userId, newPermission)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				if ok {
					continue
				}

				// Granted at user level if role is not granted
				if err := tx.Create(&models.UserRolePermission{UserId: userId, Name: newPermission, IsGranted: true}).Error; err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
			}
		}

		// Commit because permission manager tx is different
		tx.Commit()

		// Reset User Permission Cache
		permissionManager.SetUserPermissions(ctx, userId)

		return c.NoContent(http.StatusOK)
	}
}
