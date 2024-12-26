package identities

import (
	"github.com/danielgtaylor/huma/v2"
	authenticate "github.com/tguankheng016/commerce-mono/internal/identities/features/authenticating/v2"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
)

const (
	authenticateEndpoint = "/api/v2/identities/authenticate"
)

func (suite *IdentityTestSuite) TestShouldAuthenticateAsDefaultUser() {
	tests := []struct {
		username string
		password string
	}{
		{
			username: userConsts.DefaultAdminUserName,
			password: "123qwe",
		},
		{
			username: userConsts.DefaultUserUserName,
			password: "123qwe",
		},
	}

	for _, tt := range tests {
		request := authenticate.AuthenticateRequest{}
		request.UsernameOrEmailAddress = tt.username
		request.Password = tt.password

		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(request).
			SetResult(&authenticate.AuthenticateResult{}).
			Post(authenticateEndpoint)

		suite.NoError(err)

		suite.Equal(200, resp.StatusCode())

		result := resp.Result().(*authenticate.AuthenticateResult)
		suite.NotNil(result)
		suite.NotEmpty(result.AccessToken)
		suite.Greater(result.ExpireInSeconds, 0)
		suite.NotEmpty(result.RefreshToken)
		suite.Greater(result.RefreshTokenExpireInSeconds, 0)

		userId, claims, err := suite.JwtTokenHandler.ValidateToken(suite.Ctx, result.AccessToken, jwt.AccessToken)
		suite.NoError(err)

		tokenValidityKey, ok := claims[jwt.TokenValidityKey].(string)
		suite.True(ok)

		accessTokenCache, err := suite.CacheManager.Get(suite.Ctx, jwt.GenerateTokenValidityCacheKey(userId, tokenValidityKey))
		suite.NoError(err)
		suite.Equal(tokenValidityKey, accessTokenCache)

		userId, claims, err = suite.JwtTokenHandler.ValidateToken(suite.Ctx, result.RefreshToken, jwt.RefreshToken)
		suite.NoError(err)

		refreshTokenValidityKey, ok := claims[jwt.TokenValidityKey].(string)
		suite.True(ok)

		refreshTokenCache, err := suite.CacheManager.Get(suite.Ctx, jwt.GenerateTokenValidityCacheKey(userId, refreshTokenValidityKey))
		suite.NoError(err)
		suite.Equal(refreshTokenValidityKey, refreshTokenCache)
	}
}

func (suite *IdentityTestSuite) TestShouldGetErrorWhenAuthenticate() {
	tests := []struct {
		name      string
		username  string
		password  string
		wantError string
	}{
		{
			name:      "Wrong Password",
			username:  userConsts.DefaultAdminUserName,
			password:  "123123",
			wantError: "invalid username or password",
		},
		{
			name:      "Wrong Password",
			username:  userConsts.DefaultUserUserName,
			password:  "1231666",
			wantError: "invalid username or password",
		},
		{
			name:      "Empty username",
			username:  "",
			password:  "123qwe",
			wantError: "Please enter the username or email address",
		},
		{
			name:      "Empty password",
			username:  userConsts.DefaultUserUserName,
			password:  "",
			wantError: "Please enter the password",
		},
	}

	for _, tt := range tests {
		request := authenticate.AuthenticateRequest{}
		request.UsernameOrEmailAddress = tt.username
		request.Password = tt.password

		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(request).
			SetResult(&authenticate.AuthenticateResult{}).
			SetError(&huma.ErrorModel{}).
			Post(authenticateEndpoint)

		suite.NoError(err)

		suite.Equal(400, resp.StatusCode())

		result := resp.Error().(*huma.ErrorModel)
		suite.NotNil(result)

		suite.Equal(tt.wantError, result.Detail)
	}
}
