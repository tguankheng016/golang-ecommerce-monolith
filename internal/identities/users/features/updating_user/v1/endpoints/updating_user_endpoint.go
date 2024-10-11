package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func MapRoute(echo *echo.Echo, validator *validator.Validate, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/user")
	group.PUT("", updateUser(validator, permissionManager), middlewares.Authorize(permissions.PagesAdministrationUsersEdit))
}

// UpdateUser
// @Tags Users
// @Summary Update user
// @Description Update user
// @Accept json
// @Produce json
// @Param EditUserDto body EditUserDto false "EditUserDto"
// @Success 200 {object} UserDto
// @Security ApiKeyAuth
// @Router /api/v1/user [put]
func updateUser(validator *validator.Validate, permissionManager permissions.IPermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var editUserDto dtos.EditUserDto

		if err := c.Bind(&editUserDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, editUserDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := database.GetTxFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userManager := data.NewUserManager(tx)

		var user models.User
		if err := tx.First(&user, editUserDto.Id).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if err := copier.Copy(&user, &editUserDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := userManager.UpdateUser(&user, editUserDto.Password); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := userManager.UpdateUserRoles(&user, editUserDto.RoleIds); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		permissionManager.RemoveUserRoleCaches(ctx, user.Id)

		var userDto dtos.UserDto
		if err := copier.Copy(&userDto, &user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, userDto)
	}
}
