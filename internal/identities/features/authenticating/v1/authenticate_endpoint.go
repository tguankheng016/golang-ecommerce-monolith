package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tguankheng016/commerce-mono/internal/identities/services"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"github.com/tguankheng016/commerce-mono/pkg/security"
)

// Request
type AuthenticateRequest struct {
	Body struct {
		UsernameOrEmailAddress string `json:"usernameOrEmailAddress"`
		Password               string `json:"password"`
	}
}

// Result
type AuthenticateResult struct {
	Body struct {
		AccessToken                 string `json:"accessToken"`
		ExpireInSeconds             int    `json:"expireInSeconds"`
		RefreshToken                string `json:"refreshToken"`
		RefreshTokenExpireInSeconds int    `json:"refreshTokenExpireInSeconds"`
	}
}

// Validator
func (e AuthenticateRequest) Schema() v.Schema {
	return v.Schema{
		v.F("username_or_email", e.Body.UsernameOrEmailAddress): v.Nonzero[string]().Msg("Please enter the username or email address"),
		v.F("password", e.Body.Password):                        v.Nonzero[string]().Msg("Please enter the password"),
	}
}

// Handler
func MapRoute(api huma.API, jwtTokenGenerator services.IJwtTokenGenerator) {
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
		authenticate(jwtTokenGenerator),
	)
}

func authenticate(jwtTokenGenerator services.IJwtTokenGenerator) func(context.Context, *AuthenticateRequest) (*AuthenticateResult, error) {
	return func(ctx context.Context, request *AuthenticateRequest) (*AuthenticateResult, error) {
		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		userManager := userService.NewUserManager(tx)

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

		result := AuthenticateResult{}
		result.Body.AccessToken = accessToken
		result.Body.ExpireInSeconds = accessTokenSeconds
		result.Body.RefreshToken = refreshToken
		result.Body.RefreshTokenExpireInSeconds = refreshTokenSeconds

		return &result, nil
	}
}
