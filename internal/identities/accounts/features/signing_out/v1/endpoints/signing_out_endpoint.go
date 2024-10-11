package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
)

func MapRoute(echo *echo.Echo, jwt jwt.IJwtTokenValidator, jwtTokenGenerator jwt.IJwtTokenGenerator) {
	group := echo.Group("/api/v1/accounts/sign-out")
	group.POST("", signOut(jwtTokenGenerator), middlewares.ValidateToken(jwt))
}

// SigningOut
// @Tags Accounts
// @Summary Sign out
// @Description Sign out
// @Accept json
// @Produce json
// @Success 200
// @Security ApiKeyAuth
// @Router /api/v1/accounts/sign-out [post]
func signOut(jwtTokenGenerator jwt.IJwtTokenGenerator) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userId, ok := echoServer.GetCurrentUser(c)
		if ok {
			claims, ok := echoServer.GetCurrentUserClaims(c)

			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("Unable to get user claims"))
			}

			if err := jwtTokenGenerator.RemoveUserTokens(ctx, userId, claims); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}

		return c.NoContent(http.StatusOK)
	}
}
