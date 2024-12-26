package identities

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/danielgtaylor/huma/v2"
	updateCurrentProfile "github.com/tguankheng016/commerce-mono/internal/identities/features/updating_current_profile/v1"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/tests/integration_tests/users"
)

const (
	updateCurrentProfileEndpoint = "/api/v2/identities/current-profile"
)

func (suite *IdentityTestSuite) TestShouldUpdateCurrentProfileCorrectly() {
	userManager := userService.NewUserManager(suite.Pool)

	newUser := users.GetFakeUser()
	err := userManager.CreateUser(suite.Ctx, newUser, "123123")
	suite.NoError(err)

	token, err := suite.LoginAs(newUser.UserName)
	suite.NoError(err)

	request := updateCurrentProfile.UpdateCurrentProfileRequest{}
	request.FirstName = "GK66"
	request.LastName = "Tan22"
	request.Email = "testgk@666.com"
	request.UserName = "testgk666"

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetAuthToken(token).
		Put(updateCurrentProfileEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	updatedUser, err := userManager.GetUserById(suite.Ctx, newUser.Id)
	suite.NoError(err)

	suite.Equal(request.FirstName, updatedUser.FirstName)
	suite.Equal(request.LastName, updatedUser.LastName)
	suite.Equal(request.Email, updatedUser.Email)
}

func (suite *IdentityTestSuite) TestShouldUpdateCurrentProfileUnauthorized() {
	request := GetFakeUpdateCurrentProfileRequest()

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		Put(updateCurrentProfileEndpoint)

	suite.NoError(err)
	suite.Equal(401, resp.StatusCode())
}

func (suite *IdentityTestSuite) TestShouldUpdateCurrentProfileWithError() {
	userManager := userService.NewUserManager(suite.Pool)

	newUser := users.GetFakeUser()
	err := userManager.CreateUser(suite.Ctx, newUser, "123123")
	suite.NoError(err)

	tests := []struct {
		username         string
		request          *updateCurrentProfile.UpdateCurrentProfileRequest
		wantErrorMessage string
	}{
		{
			username:         newUser.UserName,
			request:          GetUpdateCurrentProfileWithErrorTestData(0),
			wantErrorMessage: "Please enter your first name",
		},
		{
			username:         newUser.UserName,
			request:          GetUpdateCurrentProfileWithErrorTestData(1),
			wantErrorMessage: "Please enter your last name",
		},
		{
			username:         newUser.UserName,
			request:          GetUpdateCurrentProfileWithErrorTestData(2),
			wantErrorMessage: "Please enter your username",
		},
		{
			username:         newUser.UserName,
			request:          GetUpdateCurrentProfileWithErrorTestData(3),
			wantErrorMessage: "Please enter your email",
		},
		{
			username:         newUser.UserName,
			request:          GetUpdateCurrentProfileWithErrorTestData(4),
			wantErrorMessage: "Please enter a valid email",
		},
		{
			username:         userConsts.DefaultAdminUserName,
			request:          GetUpdateCurrentProfileWithErrorTestData(5),
			wantErrorMessage: "You cannot update admin's username!",
		},
	}

	for _, tt := range tests {
		if tt.username == "" {
			resp, err := suite.Client.R().
				SetContext(suite.Ctx).
				SetBody(tt.request).
				Put(updateCurrentProfileEndpoint)

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
				Put(updateCurrentProfileEndpoint)

			suite.NoError(err)
			suite.Equal(400, resp.StatusCode())

			result := resp.Error().(*huma.ErrorModel)
			suite.NotNil(result)
			suite.Equal(tt.wantErrorMessage, result.Detail)
		}
	}
}

func GetUpdateCurrentProfileWithErrorTestData(scenario int) *updateCurrentProfile.UpdateCurrentProfileRequest {
	request := GetFakeUpdateCurrentProfileRequest()

	switch scenario {
	case 0:
		request.FirstName = ""
	case 1:
		request.LastName = ""
	case 2:
		request.UserName = ""
	case 3:
		request.Email = ""
	case 4:
		request.Email = "test666"
	}

	return request
}

func GetFakeUpdateCurrentProfileRequest() *updateCurrentProfile.UpdateCurrentProfileRequest {
	updateCurrentProfile := updateCurrentProfile.UpdateCurrentProfileRequest{}
	updateCurrentProfile.FirstName = gofakeit.FirstName()
	updateCurrentProfile.LastName = gofakeit.LastName()
	updateCurrentProfile.UserName = gofakeit.Username()
	updateCurrentProfile.Email = gofakeit.Email()

	return &updateCurrentProfile
}
