package endpoints

import (
	"context"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/data"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func MapRoute(jwt jwt.IJwtTokenValidator, checker permissions.IPermissionChecker, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/user")
	group.POST("", createUser(log, ctx), middlewares.ValidateToken(jwt), middlewares.Authorize(checker, permissions.PagesAdministrationUsersCreate))
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
func createUser(log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		tx, err := database.RetrieveTxContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userManager := data.NewUserManager(tx)

		var createUserDto dtos.CreateUserDto

		if err := c.Bind(&createUserDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		var user models.User
		copier.Copy(&user, &createUserDto)

		if err := userManager.CreateUser(&user, createUserDto.Password); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var userDto dtos.UserDto
		copier.Copy(&userDto, &user)

		return c.JSON(http.StatusOK, userDto)
	}
}
