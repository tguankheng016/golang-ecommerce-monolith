package endpoints

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/accounts/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

type GetCurrentSessionResult struct {
	User               *dtos.UserLoginInfoDto `json:"user"`
	AllPermissions     map[string]bool        `json:"allPermissions"`
	GrantedPermissions map[string]bool        `json:"grantedPermissions"`
} // @name GetCurrentSessionResult

func MapRoute(echo *echo.Echo, jwt jwt.IJwtTokenValidator, permissionManager permissions.IPermissionManager) {
	group := echo.Group("/api/v1/accounts/current-session")
	group.GET("", getCurrentSession(permissionManager), middlewares.TryValidateToken(jwt))
}

// GetCurrentSession
// @Tags Accounts
// @Summary Get Current User Session
// @Description Get Current User Session
// @Accept json
// @Produce json
// @Success 200 {object} GetCurrentSessionResult
// @Security ApiKeyAuth
// @Router /api/v1/accounts/current-session [get]
func getCurrentSession(permissionManager permissions.IPermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		allPermisions := permissions.GetAppPermissions().Items
		currentUserSessionResult := GetCurrentSessionResult{
			AllPermissions: convertPermissionMap(allPermisions),
		}

		userId, ok := echoServer.GetCurrentUser(c)
		if ok {
			tx, err := database.RetrieveTxCtx(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			currentUserSessionResult.User = &dtos.UserLoginInfoDto{}
			var user models.User
			if err := tx.First(&user, userId).Error; err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			if err := copier.Copy(&currentUserSessionResult.User, &user); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			userGrantedPermissions, err := permissionManager.SetUserPermissions(ctx, userId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			currentUserSessionResult.GrantedPermissions = convertPermissionMap(userGrantedPermissions)

		} else {
			currentUserSessionResult.GrantedPermissions = make(map[string]bool)
		}

		return c.JSON(http.StatusOK, currentUserSessionResult)
	}
}

func convertPermissionMap[V any](permissionMap map[string]V) map[string]bool {
	boolMap := make(map[string]bool)
	for key := range permissionMap {
		boolMap[key] = true
	}
	return boolMap
}
