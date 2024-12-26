package identities

import (
	"github.com/danielgtaylor/huma/v2"
	authenticate "github.com/tguankheng016/commerce-mono/internal/identities/features/authenticating/v2"
	changePassword "github.com/tguankheng016/commerce-mono/internal/identities/features/changing_password/v1"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
)

const (
	changePasswordEndpoint = "/api/v2/identities/change-password"
)

func (suite *IdentityTestSuite) TestShouldChangePasswordCorrectly() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	request := changePassword.ChangePasswordRequest{}
	request.CurrentPassword = "123qwe"
	request.NewPassword = "123123123"

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetAuthToken(token).
		Put(changePasswordEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	authRequest := authenticate.AuthenticateRequest{}
	authRequest.UsernameOrEmailAddress = userConsts.DefaultUserUserName
	authRequest.Password = request.NewPassword

	resp, err = suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(authRequest).
		SetResult(&authenticate.AuthenticateResult{}).
		Post(authenticateEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*authenticate.AuthenticateResult)
	suite.NotNil(result)
}

func (suite *IdentityTestSuite) TestShouldChangePasswordWithError() {
	tests := []struct {
		username         string
		request          *changePassword.ChangePasswordRequest
		wantErrorMessage string
	}{
		{
			username: userConsts.DefaultUserUserName,
			request: &changePassword.ChangePasswordRequest{
				CurrentPassword: "1231231",
				NewPassword:     "123123123",
			},
			wantErrorMessage: "current password is incorrect",
		},
		{
			username: "",
			request: &changePassword.ChangePasswordRequest{
				CurrentPassword: "123qwe",
				NewPassword:     "123123123",
			},
			wantErrorMessage: "",
		},
	}

	for _, tt := range tests {
		if tt.username == "" {
			resp, err := suite.Client.R().
				SetContext(suite.Ctx).
				SetBody(tt.request).
				Put(changePasswordEndpoint)

			suite.NoError(err)
			suite.Equal(401, resp.StatusCode())
		} else {
			token, err := suite.LoginAs(tt.username)
			suite.NoError(err)

			resp, err := suite.Client.R().
				SetContext(suite.Ctx).
				SetBody(tt.request).
				SetError(&huma.ErrorModel{}).
				SetAuthToken(token).
				Put(changePasswordEndpoint)

			suite.NoError(err)
			suite.Equal(400, resp.StatusCode())

			result := resp.Error().(*huma.ErrorModel)
			suite.NotNil(result)
			suite.Equal(tt.wantErrorMessage, result.Detail)
		}
	}
}

func (suite *IdentityTestSuite) TestShouldChangePasswordUnauthorized() {
	request := changePassword.ChangePasswordRequest{}
	request.CurrentPassword = "123qwe"
	request.NewPassword = "123123123"

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		Put(changePasswordEndpoint)

	suite.NoError(err)
	suite.Equal(401, resp.StatusCode())
}
