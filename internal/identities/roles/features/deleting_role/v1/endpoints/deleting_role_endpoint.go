package endpoints

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func MapRoute(echo *echo.Echo, jwt jwt.IJwtTokenValidator, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/role/:roleId")
	group.DELETE("", deleteRole(permissionManager), middlewares.ValidateToken(jwt), middlewares.Authorize(permissionManager, permissions.PagesAdministrationRolesDelete))
}

// DeleteRole
// @Tags Roles
// @Summary Delete role
// @Description Delete role
// @Accept json
// @Produce json
// @Param roleId path int true "Role Id"
// @Success 200
// @Security ApiKeyAuth
// @Router /api/v1/role/{roleId} [delete]
func deleteRole(permissionManager permissions.IPermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var roleId int64

		err := echo.PathParamsBinder(c).
			Int64("roleId", &roleId).
			BindError()

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := database.RetrieveTxCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var role models.Role
		if err := tx.First(&role, roleId).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if strings.EqualFold(role.Name, constants.DefaultAdminRoleName) {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("You cannot delete the default admin user"))
		}

		userManager := data.NewUserManager(tx)
		userIds, err := userManager.GetUserIdsInRole(roleId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		for _, userId := range userIds {
			if err := userManager.RemoveToRoles(&models.User{Id: userId}, []int64{roleId}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			permissionManager.RemoveUserRoleCaches(ctx, userId)
		}

		if err := tx.Delete(&role).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusOK)
	}
}
