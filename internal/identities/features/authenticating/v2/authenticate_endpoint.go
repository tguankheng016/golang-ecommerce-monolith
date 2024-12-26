package v2

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/identities/services"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/security"
)

// Request
type AuthenticateRequest struct {
	UsernameOrEmailAddress string `json:"usernameOrEmailAddress"`
	Password               string `json:"password"`
}
type HumaAuthenticateRequest struct {
	Body struct {
		AuthenticateRequest
	}
}

// Result
type AuthenticateResult struct {
	AccessToken                 string `json:"accessToken"`
	ExpireInSeconds             int    `json:"expireInSeconds"`
	RefreshToken                string `json:"refreshToken"`
	RefreshTokenExpireInSeconds int    `json:"refreshTokenExpireInSeconds"`
}
type HumaAuthenticateResult struct {
	Body struct {
		AuthenticateResult
	}
}

// Validator
func (e HumaAuthenticateRequest) Schema() v.Schema {
	return v.Schema{
		v.F("username_or_email", e.Body.UsernameOrEmailAddress): v.Nonzero[string]().Msg("Please enter the username or email address"),
		v.F("password", e.Body.Password):                        v.Nonzero[string]().Msg("Please enter the password"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	jwtTokenGenerator services.IJwtTokenGenerator,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "Authenticate",
			Method:        http.MethodPost,
			Path:          "/identities/authenticate",
			Summary:       "Authenticate",
			Tags:          []string{"Identities"},
			DefaultStatus: http.StatusOK,
		},
		authenticate(jwtTokenGenerator, pool),
	)
}

func authenticate(jwtTokenGenerator services.IJwtTokenGenerator, pool *pgxpool.Pool) func(context.Context, *HumaAuthenticateRequest) (*HumaAuthenticateResult, error) {
	return func(ctx context.Context, request *HumaAuthenticateRequest) (*HumaAuthenticateResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		userManager := userService.NewUserManager(pool)

		user, err := userManager.GetUserByUserName(ctx, request.Body.UsernameOrEmailAddress)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error(), err)
		}

		if user == nil {
			user, err := userManager.GetUserByEmail(ctx, request.Body.UsernameOrEmailAddress)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error(), err)
			}

			if user == nil {
				return nil, huma.Error404NotFound("user not found")
			}
		}

		ok, err := security.ComparePasswords(user.PasswordHash, request.Body.Password)
		if err != nil || !ok {
			return nil, huma.Error400BadRequest("invalid username or password")
		}

		refreshToken, refreshTokenKey, refreshTokenSeconds, err := jwtTokenGenerator.GenerateRefreshToken(ctx, *user)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error(), err)
		}

		accessToken, accessTokenSeconds, err := jwtTokenGenerator.GenerateAccessToken(ctx, *user, refreshTokenKey)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error(), err)
		}

		result := HumaAuthenticateResult{}

		result.Body.AccessToken = accessToken
		result.Body.ExpireInSeconds = accessTokenSeconds
		result.Body.RefreshToken = refreshToken
		result.Body.RefreshTokenExpireInSeconds = refreshTokenSeconds

		return &result, nil
	}
}
