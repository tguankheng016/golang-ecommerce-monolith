package endpoints

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

func MapRoute(echo *echo.Echo, jwt jwt.IJwtTokenValidator, checker permissions.IPermissionChecker) {
	group := echo.Group("/api/v1/user/:userId")
	group.DELETE("", deleteUser(), middlewares.ValidateToken(jwt), middlewares.Authorize(checker, permissions.PagesAdministrationUsersDelete))
}

// DeleteUser
// @Tags Users
// @Summary Delete user
// @Description Delete user
// @Accept json
// @Produce json
// @Param userId path int true "User Id"
// @Success 200
// @Security ApiKeyAuth
// @Router /api/v1/user/{userId} [delete]
func deleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var userId int64

		err := echo.PathParamsBinder(c).
			Int64("userId", &userId).
			BindError()

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := database.RetrieveTxCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var user models.User
		if err := tx.First(&user, userId).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if user.NormalizedUserName == strings.ToUpper(constants.DefaultAdminUsername) {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("You cannot delete the default admin user"))
		}

		if err := tx.Delete(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusOK)
	}
}
