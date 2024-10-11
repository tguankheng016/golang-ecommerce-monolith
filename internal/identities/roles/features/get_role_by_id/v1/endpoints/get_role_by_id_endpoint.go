package endpoints

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/roles/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/core/helpers"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

type GetRoleByIdResult struct {
	Role dtos.CreateOrEditRoleDto
} // @name GetRoleByIdResult

func MapRoute(echo *echo.Echo, jwt jwt.IJwtTokenValidator, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/role/:roleId")
	group.GET("", getRoleById(permissionManager), middlewares.ValidateToken(jwt), middlewares.Authorize(permissionManager, permissions.PagesAdministrationRoles))
}

// GetRoleById
// @Tags Roles
// @Summary Get role by id
// @Description Get role by id
// @Accept json
// @Produce json
// @Param roleId path int true "Role Id"
// @Success 200 {object} GetRoleByIdResult
// @Security ApiKeyAuth
// @Router /api/v1/role/{roleId} [get]
func getRoleById(permissionManager permissions.IPermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var roleId int64

		err := echo.PathParamsBinder(c).
			Int64("roleId", &roleId).
			BindError()

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var roleEditDto dtos.CreateOrEditRoleDto

		if roleId == 0 {
			// Create
			roleEditDto = dtos.CreateOrEditRoleDto{}
		} else {
			// Edit
			tx, err := database.RetrieveTxCtx(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			var role models.Role
			if err := tx.First(&role, roleId).Error; err != nil {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}

			if err := copier.Copy(&roleEditDto, &role); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			rolePermissions, err := permissionManager.SetRolePermissions(ctx, role.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			roleEditDto.GrantedPermissions = helpers.MapKeysToSlice(rolePermissions)
		}

		result := &GetRoleByIdResult{
			Role: roleEditDto,
		}

		return c.JSON(http.StatusOK, result)
	}
}
