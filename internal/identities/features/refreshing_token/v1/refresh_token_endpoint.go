package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/identities/services"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
)

// Request
type RefreshTokenRequest struct {
	Token string `json:"token"`
}
type HumaRefreshTokenRequest struct {
	Body struct {
		RefreshTokenRequest
	}
}

// Result
type RefreshTokenResult struct {
	AccessToken     string `json:"accessToken"`
	ExpireInSeconds int    `json:"expireInSeconds"`
}
type HumaRefreshTokenResult struct {
	Body struct {
		RefreshTokenResult
	}
}

// Validator
func (e HumaRefreshTokenRequest) Schema() v.Schema {
	return v.Schema{
		v.F("token", e.Body.Token): v.Nonzero[string]().Msg("Refresh token cannot be empty"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	jwtTokenGenerator services.IJwtTokenGenerator,
	jwtTokenHandler jwt.IJwtTokenHandler,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "RefreshToken",
			Method:        http.MethodPost,
			Path:          "/identities/refresh-token",
			Summary:       "RefreshToken",
			Tags:          []string{"Identities"},
			DefaultStatus: http.StatusOK,
		},
		refreshToken(pool, jwtTokenGenerator, jwtTokenHandler),
	)
}

func refreshToken(pool *pgxpool.Pool, jwtTokenGenerator services.IJwtTokenGenerator, jwtTokenHandler jwt.IJwtTokenHandler) func(context.Context, *HumaRefreshTokenRequest) (*HumaRefreshTokenResult, error) {
	return func(ctx context.Context, request *HumaRefreshTokenRequest) (*HumaRefreshTokenResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		userId, claims, err := jwtTokenHandler.ValidateToken(ctx, request.Body.Token, jwt.RefreshToken)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		userManager := userService.NewUserManager(pool)

		user, err := userManager.GetUserById(ctx, userId)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if user == nil {
			return nil, huma.Error404NotFound("user not found")
		}

		refreshTokenKey := claims[jwt.TokenValidityKey].(string)
		accessToken, accessTokenSeconds, err := jwtTokenGenerator.GenerateAccessToken(ctx, *user, refreshTokenKey)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := HumaRefreshTokenResult{}
		result.Body.AccessToken = accessToken
		result.Body.ExpireInSeconds = accessTokenSeconds

		return &result, nil
	}
}
