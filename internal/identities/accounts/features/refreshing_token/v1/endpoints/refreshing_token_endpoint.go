package endpoints

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	appConstants "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"gorm.io/gorm"
)

type RefreshTokenRequest struct {
	Token string `json:"token"`
} // @name RefreshTokenRequest

type RefreshTokenResult struct {
	AccessToken     string `json:"accessToken"`
	ExpireInSeconds int    `json:"expireInSeconds"`
} // @name RefreshTokenResult

func MapRoute(db *gorm.DB, jwtTokenGenerator jwt.IJwtTokenGenerator, jwtTokenValidator jwt.IJwtTokenValidator, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/accounts/refresh-token")
	group.POST("", refreshToken(db, jwtTokenGenerator, jwtTokenValidator, log, ctx))
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
func refreshToken(db *gorm.DB, jwtTokenGenerator jwt.IJwtTokenGenerator, jwtTokenValidator jwt.IJwtTokenValidator, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &RefreshTokenRequest{}

		if err := c.Bind(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		claims, err := jwtTokenValidator.ValidateToken(request.Token, jwt.RefreshToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("Invalid token"))
		}

		userId, err := strconv.ParseInt(sub, 10, 64)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var user models.User

		if err := db.First(&user, userId).Error; err != nil {
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
