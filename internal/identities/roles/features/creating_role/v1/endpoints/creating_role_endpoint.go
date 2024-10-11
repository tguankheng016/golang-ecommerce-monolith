package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/roles/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func MapRoute(echo *echo.Echo, validator *validator.Validate) {
	group := echo.Group("/api/v1/role")
	group.POST("", createRole(validator), middlewares.Authorize(permissions.PagesAdministrationRolesCreate))
}

// CreateRole
// @Tags Roles
// @Summary Create new role
// @Description Create new role
// @Accept json
// @Produce json
// @Param CreateRoleDto body CreateRoleDto false "CreateRoleDto"
// @Success 200 {object} RoleDto
// @Security ApiKeyAuth
// @Router /api/v1/role [post]
func createRole(validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := database.GetTxFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		roleManager := data.NewRoleManager(tx)

		var createRoleDto dtos.CreateRoleDto

		if err := c.Bind(&createRoleDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, createRoleDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var role models.Role
		if err := copier.Copy(&role, &createRoleDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := permissions.ValidatePermissionName(createRoleDto.GrantedPermissions); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := roleManager.CreateRole(&role); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if len(createRoleDto.GrantedPermissions) > 0 {
			for _, permission := range createRoleDto.GrantedPermissions {
				newUserRolePermission := &models.UserRolePermission{
					RoleId:    role.Id,
					Name:      permission,
					IsGranted: true,
				}
				if err := tx.Model(&models.UserRolePermission{}).Create(&newUserRolePermission).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
			}
		}

		var roleDto dtos.RoleDto
		if err := copier.Copy(&roleDto, &role); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, roleDto)
	}
}
