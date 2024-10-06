package endpoints

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/security"
	"gorm.io/gorm"
)

type AuthenticateRequest struct {
	UsernameOrEmailAddress string `json:"usernameOrEmailAddress"`
	Password               string `json:"password"`
} // @name AuthenticateRequest

type AuthenticateResult struct {
	AccessToken                 string `json:"accessToken"`
	ExpireInSeconds             int    `json:"expireInSeconds"`
	RefreshToken                string `json:"refreshToken"`
	RefreshTokenExpireInSeconds int    `json:"refreshTokenExpireInSeconds"`
} // @name AuthenticateResult

func MapRoute(db *gorm.DB, jwtTokenGenerator jwt.IJwtTokenGenerator, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/accounts/authenticate")
	group.POST("", authenticate(db, jwtTokenGenerator, log, ctx))
}

// Authenticate
// @Tags Accounts
// @Summary Authenticate
// @Description Authenticate
// @Accept json
// @Produce json
// @Param AuthenticateRequest body AuthenticateRequest true "AuthenticateRequest"
// @Success 200 {object} AuthenticateResult
// @Security ApiKeyAuth
// @Router /api/v1/accounts/authenticate [post]
func authenticate(db *gorm.DB, jwtTokenGenerator jwt.IJwtTokenGenerator, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &AuthenticateRequest{}

		if err := c.Bind(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var user models.User

		err := db.Where("user_name = ? OR email = ?", request.UsernameOrEmailAddress, request.UsernameOrEmailAddress).First(&user).Error
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		ok, err := security.ComparePasswords(user.Password, request.Password)
		if err != nil || !ok {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		refreshToken, refreshTokenKey, err := jwtTokenGenerator.GenerateRefreshToken(&user)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		accessToken, err := jwtTokenGenerator.GenerateAccessToken(&user, refreshTokenKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result := &AuthenticateResult{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return c.JSON(http.StatusOK, result)
	}
}
