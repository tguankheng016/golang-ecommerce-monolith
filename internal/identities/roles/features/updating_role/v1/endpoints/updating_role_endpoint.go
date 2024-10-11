package endpoints

import (
	"net/http"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/roles/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
	"gorm.io/gorm"
)

func MapRoute(echo *echo.Echo, validator *validator.Validate, jwt jwt.IJwtTokenValidator, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/role")
	group.PUT("", updateRole(validator, permissionManager), middlewares.ValidateToken(jwt), middlewares.Authorize(permissionManager, permissions.PagesAdministrationRolesEdit))
}

// UpdateRole
// @Tags Roles
// @Summary Update role
// @Description Update role
// @Accept json
// @Produce json
// @Param EditRoleDto body EditRoleDto false "EditRoleDto"
// @Success 200 {object} RoleDto
// @Security ApiKeyAuth
// @Router /api/v1/role [put]
func updateRole(validator *validator.Validate, permissionManager permissions.IPermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var editRoleDto dtos.EditRoleDto

		if err := c.Bind(&editRoleDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, editRoleDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := database.RetrieveTxCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		roleManager := data.NewRoleManager(tx)

		var role models.Role
		if err := tx.First(&role, editRoleDto.Id).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if err := copier.Copy(&role, &editRoleDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := permissions.ValidatePermissionName(editRoleDto.GrantedPermissions); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := roleManager.UpdateRole(&role); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		oldPermissions, err := permissionManager.SetRolePermissions(ctx, role.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		isAdmin := strings.EqualFold(role.Name, constants.DefaultAdminRoleName)

		// Prohibit
		for _, oldPermission := range oldPermissions {
			if !slices.Contains(editRoleDto.GrantedPermissions, oldPermission) {
				if err := tx.Where("role_id = ? AND name = ?", role.Id, oldPermission).Delete(&models.UserRolePermission{}).Error; err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}

				if isAdmin {
					if err := tx.Create(&models.UserRolePermission{RoleId: role.Id, Name: oldPermission, IsGranted: false}).Error; err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err)
					}
				}
			}
		}

		// Granted
		for _, newPermission := range editRoleDto.GrantedPermissions {
			if _, ok := oldPermissions[newPermission]; !ok {
				var rolePermissionToGrant models.UserRolePermission
				if err := tx.First(&rolePermissionToGrant, models.UserRolePermission{RoleId: role.Id, Name: newPermission}).Error; err != nil && err != gorm.ErrRecordNotFound {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				if rolePermissionToGrant.Id == 0 {
					if err := tx.Create(&models.UserRolePermission{RoleId: role.Id, Name: newPermission, IsGranted: true}).Error; err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err)
					}
				} else if !rolePermissionToGrant.IsGranted {
					if isAdmin {
						if err := tx.Where("role_id = ? AND name = ?", role.Id, newPermission).Delete(&models.UserRolePermission{}).Error; err != nil {
							return echo.NewHTTPError(http.StatusInternalServerError, err)
						}
					} else {
						// Unlikely will hit this case
						rolePermissionToGrant.IsGranted = true
						if err := tx.Save(&rolePermissionToGrant).Error; err != nil {
							return echo.NewHTTPError(http.StatusInternalServerError, err)
						}
					}
				}
			}
		}

		permissionManager.RemoveRolePermissionCaches(ctx, role.Id)

		var roleDto dtos.RoleDto
		if err := copier.Copy(&roleDto, &role); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, roleDto)
	}
}
