package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/security"
)

type AuthenticateRequest struct {
	UsernameOrEmailAddress string `json:"usernameOrEmailAddress" validate:"required"`
	Password               string `json:"password" validate:"required"`
} // @name AuthenticateRequest

type AuthenticateResult struct {
	AccessToken                 string `json:"accessToken"`
	ExpireInSeconds             int    `json:"expireInSeconds"`
	RefreshToken                string `json:"refreshToken"`
	RefreshTokenExpireInSeconds int    `json:"refreshTokenExpireInSeconds"`
} // @name AuthenticateResult

func MapRoute(echo *echo.Echo, validator *validator.Validate, jwtTokenGenerator jwt.IJwtTokenGenerator) {
	group := echo.Group("/api/v1/accounts/authenticate")
	group.POST("", authenticate(validator, jwtTokenGenerator))
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
func authenticate(validator *validator.Validate, jwtTokenGenerator jwt.IJwtTokenGenerator) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &AuthenticateRequest{}

		if err := c.Bind(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := database.GetTxFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var user models.User
		if err := tx.Where("user_name = ? OR email = ?", request.UsernameOrEmailAddress, request.UsernameOrEmailAddress).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		ok, err := security.ComparePasswords(user.Password, request.Password)
		if err != nil || !ok {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		refreshToken, refreshTokenKey, refreshTokenSeconds, err := jwtTokenGenerator.GenerateRefreshToken(ctx, &user)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		accessToken, accessTokenSeconds, err := jwtTokenGenerator.GenerateAccessToken(ctx, &user, refreshTokenKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result := &AuthenticateResult{
			AccessToken:                 accessToken,
			ExpireInSeconds:             accessTokenSeconds,
			RefreshToken:                refreshToken,
			RefreshTokenExpireInSeconds: refreshTokenSeconds,
		}

		return c.JSON(http.StatusOK, result)
	}
}
