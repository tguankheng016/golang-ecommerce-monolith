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
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func MapRoute(echo *echo.Echo, validator *validator.Validate, jwt jwt.IJwtTokenValidator, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/user")
	group.POST("", createUser(validator), middlewares.ValidateToken(jwt), middlewares.Authorize(permissionManager, permissions.PagesAdministrationUsersCreate))
}

// CreateUser
// @Tags Users
// @Summary Create new user
// @Description Create new user
// @Accept json
// @Produce json
// @Param CreateUserDto body CreateUserDto false "CreateUserDto"
// @Success 200 {object} UserDto
// @Security ApiKeyAuth
// @Router /api/v1/user [post]
func createUser(validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := database.RetrieveTxCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userManager := data.NewUserManager(tx)

		var createUserDto dtos.CreateUserDto

		if err := c.Bind(&createUserDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, createUserDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var user models.User
		if err := copier.Copy(&user, &createUserDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := userManager.CreateUser(&user, createUserDto.Password); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var userDto dtos.UserDto
		if err := copier.Copy(&userDto, &user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, userDto)
	}
}
