package endpoints

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

type GetUserByIdResult struct {
	User dtos.CreateOrEditUserDto
} // @name GetUserByIdResult

func MapRoute(echo *echo.Echo, jwt jwt.IJwtTokenValidator, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/user/:userId")
	group.GET("", getUserById(), middlewares.ValidateToken(jwt), middlewares.Authorize(permissionManager, permissions.PagesAdministrationUsers))
}

// GetUserById
// @Tags Users
// @Summary Get user by id
// @Description Get user by id
// @Accept json
// @Produce json
// @Param userId path int true "User Id"
// @Success 200 {object} GetUserByIdResult
// @Security ApiKeyAuth
// @Router /api/v1/user/{userId} [get]
func getUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var userId int64
		if err := echo.PathParamsBinder(c).Int64("userId", &userId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var userEditDto dtos.CreateOrEditUserDto

		if userId == 0 {
			// Create
			userEditDto = dtos.CreateOrEditUserDto{}
		} else {
			// Edit
			tx, err := database.RetrieveTxCtx(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			var user models.User
			if err := tx.First(&user, userId).Error; err != nil {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}

			if err := copier.Copy(&userEditDto, &user); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			userManager := data.NewUserManager(tx)

			userEditDto.RoleIds, err = userManager.GetUserRoles(&user)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}

		result := &GetUserByIdResult{
			User: userEditDto,
		}

		return c.JSON(http.StatusOK, result)
	}
}
