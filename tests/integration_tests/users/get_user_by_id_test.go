package users

import (
	"fmt"

	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/internal/users/dtos"
)

const (
	getUserByIdEndpoint = "/api/v1/user"
)

type GetUserByIdResult struct {
	User dtos.CreateOrEditUserDto
}

func (suite *UserTestSuite) TestShouldGetCorrectUserById() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		id           int64
		wantUsername string
	}{
		{
			id:           1,
			wantUsername: userConsts.DefaultAdminUserName,
		},
		{
			id:           0,
			wantUsername: "",
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetResult(&GetUserByIdResult{}).
			SetAuthToken(token).
			Get(fmt.Sprintf("%s/%d", getUserByIdEndpoint, tt.id))

		suite.NoError(err)
		suite.Equal(200, resp.StatusCode())

		result := resp.Result().(*GetUserByIdResult)
		suite.NotNil(result)
		suite.NotNil(result.User)
		suite.Equal(result.User.UserName, tt.wantUsername)

		if tt.id == 0 {
			suite.Nil(result.User.Id)
		} else {
			suite.Equal(int64(1), *result.User.Id)
		}
	}
}

func (suite *UserTestSuite) TestShouldGetErrorWhenGetUserById() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		id            int64
		wantErrorCode int
	}{
		{
			id:            -3,
			wantErrorCode: 400,
		},
		{
			id:            100,
			wantErrorCode: 404,
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetResult(&GetUserByIdResult{}).
			SetAuthToken(token).
			Get(fmt.Sprintf("%s/%d", getUserByIdEndpoint, tt.id))

		suite.NoError(err)
		suite.Equal(tt.wantErrorCode, resp.StatusCode())
	}
}

func (suite *UserTestSuite) TestShouldGetUserByIdUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetResult(&GetUserByIdResult{}).
		SetAuthToken(token).
		Get(fmt.Sprintf("%s/%d", getUserByIdEndpoint, 1))

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}
