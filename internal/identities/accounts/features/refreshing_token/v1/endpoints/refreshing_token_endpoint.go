package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	appConstants "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
)

type RefreshTokenRequest struct {
	Token string `json:"token" validate:"min=10,required"`
} // @name RefreshTokenRequest

type RefreshTokenResult struct {
	AccessToken     string `json:"accessToken"`
	ExpireInSeconds int    `json:"expireInSeconds"`
} // @name RefreshTokenResult

func MapRoute(echo *echo.Echo, validator *validator.Validate, jwtTokenGenerator jwt.IJwtTokenGenerator, jwtTokenValidator jwt.IJwtTokenValidator) {
	group := echo.Group("/api/v1/accounts/refresh-token")
	group.POST("", refreshToken(validator, jwtTokenGenerator, jwtTokenValidator))
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
func refreshToken(validator *validator.Validate, jwtTokenGenerator jwt.IJwtTokenGenerator, jwtTokenValidator jwt.IJwtTokenValidator) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &RefreshTokenRequest{}

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

		userId, claims, err := jwtTokenValidator.ValidateToken(ctx, request.Token, jwt.RefreshToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		var user models.User
		if err := tx.First(&user, userId).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		refreshTokenKey := claims[appConstants.TokenValidityKey].(string)

		accessToken, accessTokenSeconds, err := jwtTokenGenerator.GenerateAccessToken(ctx, &user, refreshTokenKey)
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
