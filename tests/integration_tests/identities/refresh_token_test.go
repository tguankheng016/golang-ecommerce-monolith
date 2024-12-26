package identities

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	authenticate "github.com/tguankheng016/commerce-mono/internal/identities/features/authenticating/v2"
	refreshToken "github.com/tguankheng016/commerce-mono/internal/identities/features/refreshing_token/v1"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
)

const (
	refreshTokenEndpoint = "/api/v2/identities/refresh-token"
)

func (suite *IdentityTestSuite) TestShouldRefreshToken() {
	authenticateRequest := authenticate.AuthenticateRequest{}
	authenticateRequest.UsernameOrEmailAddress = userConsts.DefaultAdminUserName
	authenticateRequest.Password = "123qwe"

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(authenticateRequest).
		SetResult(&authenticate.AuthenticateResult{}).
		Post(authenticateEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	authenticateResult := resp.Result().(*authenticate.AuthenticateResult)
	suite.NotNil(authenticateResult)

	request := refreshToken.RefreshTokenRequest{}
	request.Token = authenticateResult.RefreshToken

	resp, err = suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&refreshToken.RefreshTokenResult{}).
		Post(refreshTokenEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*refreshToken.RefreshTokenResult)
	suite.NotNil(result)
	suite.NotEmpty(result.AccessToken)
	suite.Greater(result.ExpireInSeconds, 0)

	userId, claims, err := suite.JwtTokenHandler.ValidateToken(suite.Ctx, result.AccessToken, jwt.AccessToken)
	suite.NoError(err)

	tokenValidityKey, ok := claims[jwt.TokenValidityKey].(string)
	suite.True(ok)

	accessTokenCache, err := suite.CacheManager.Get(suite.Ctx, jwt.GenerateTokenValidityCacheKey(userId, tokenValidityKey))
	suite.NoError(err)
	suite.Equal(tokenValidityKey, accessTokenCache)
}

func (suite *IdentityTestSuite) TestShouldExpiredRefreshTokenError() {
	authenticateRequest := authenticate.AuthenticateRequest{}
	authenticateRequest.UsernameOrEmailAddress = userConsts.DefaultAdminUserName
	authenticateRequest.Password = "123qwe"

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(authenticateRequest).
		SetResult(&authenticate.AuthenticateResult{}).
		Post(authenticateEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	authenticateResult := resp.Result().(*authenticate.AuthenticateResult)
	suite.NotNil(authenticateResult)

	request := refreshToken.RefreshTokenRequest{}
	request.Token = authenticateResult.RefreshToken

	userId, claims, err := suite.JwtTokenHandler.ValidateToken(suite.Ctx, authenticateResult.RefreshToken, jwt.RefreshToken)
	suite.NoError(err)

	tokenValidityKey, ok := claims[jwt.TokenValidityKey].(string)
	suite.True(ok)

	query := "UPDATE user_tokens set expiration_time = $1 WHERE user_id = $2 and token_key = $3"
	_, err = suite.Pool.Exec(suite.Ctx, query, time.Now().Add(-1*time.Hour), userId, tokenValidityKey)
	suite.NoError(err)

	err = suite.CacheManager.Clear(suite.Ctx)
	suite.NoError(err)

	resp, err = suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&refreshToken.RefreshTokenResult{}).
		SetError(&huma.ErrorModel{}).
		Post(refreshTokenEndpoint)

	suite.NoError(err)
	suite.Equal(500, resp.StatusCode())

	result := resp.Error().(*huma.ErrorModel)

	suite.Equal("invalid token key", result.Detail)
}

func (suite *IdentityTestSuite) TestShouldGetErrorWhenRefreshToken() {
	tests := []struct {
		token            string
		wantErrorCode    int
		wantErrorMessage string
	}{
		{
			token:            "",
			wantErrorCode:    400,
			wantErrorMessage: "Refresh token cannot be empty",
		},
		{
			token:            "invalid-token",
			wantErrorCode:    500,
			wantErrorMessage: "",
		},
	}

	for _, tt := range tests {
		request := refreshToken.RefreshTokenRequest{}
		request.Token = tt.token

		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(request).
			SetResult(&refreshToken.RefreshTokenResult{}).
			SetError(&huma.ErrorModel{}).
			Post(refreshTokenEndpoint)

		suite.NoError(err)
		suite.Equal(tt.wantErrorCode, resp.StatusCode())

		result := resp.Error().(*huma.ErrorModel)
		suite.NotNil(result)

		if tt.wantErrorMessage != "" {
			suite.Equal(tt.wantErrorMessage, result.Detail)
		}
	}
}
