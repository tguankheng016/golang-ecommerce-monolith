package endpoints

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	appConstants "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
)

type RefreshTokenRequest struct {
	Token string `json:"token"`
} // @name RefreshTokenRequest

type RefreshTokenResult struct {
	AccessToken     string `json:"accessToken"`
	ExpireInSeconds int    `json:"expireInSeconds"`
} // @name RefreshTokenResult

func MapRoute(jwtTokenGenerator jwt.IJwtTokenGenerator, jwtTokenValidator jwt.IJwtTokenValidator, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/accounts/refresh-token")
	group.POST("", refreshToken(jwtTokenGenerator, jwtTokenValidator, log, ctx))
}

// RefreshToken
// @Tags Accounts
// @Summary Refresh access token
// @Description Refresh access token
// @Accept json
// @Produce json
// @Param RefreshTokenRequest body RefreshTokenRequest true "RefreshTokenRequest"
// @Success 200 {object} RefreshTokenResult
// @Security ApiKeyAuth
// @Router /api/v1/accounts/refresh-token [post]
func refreshToken(jwtTokenGenerator jwt.IJwtTokenGenerator, jwtTokenValidator jwt.IJwtTokenValidator, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &RefreshTokenRequest{}

		if err := c.Bind(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := database.RetrieveTxContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userId, claims, err := jwtTokenValidator.ValidateToken(request.Token, jwt.RefreshToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		var user models.User
		if err := tx.First(&user, userId).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		refreshTokenKey := claims[appConstants.TokenValidityKey].(string)

		accessToken, accessTokenSeconds, err := jwtTokenGenerator.GenerateAccessToken(&user, refreshTokenKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result := &RefreshTokenResult{
			AccessToken:     accessToken,
			ExpireInSeconds: accessTokenSeconds,
		}

		return c.JSON(http.StatusOK, result)
	}
}
