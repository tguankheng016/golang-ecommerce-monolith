package users

import (
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/internal/users/dtos"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/core/pagination"
)

const (
	getUsersEndpoint = "/api/v2/users"
)

func (suite *UserTestSuite) TestShouldGetUsers() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	userManager := userService.NewUserManager(suite.Pool)

	totalCount, err := userManager.GetUsersCount(suite.Ctx)
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetResult(&pagination.PageResultDto[dtos.UserDto]{}).
		SetAuthToken(token).
		Get(getUsersEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*pagination.PageResultDto[dtos.UserDto])
	suite.NotNil(result)
	suite.Equal(totalCount, result.TotalCount)
	suite.Equal(totalCount, len(result.Items))
}

func (suite *UserTestSuite) TestShouldGetFilteredUsers() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetQueryParams(map[string]string{
			"filters": "admi",
		}).
		SetResult(&pagination.PageResultDto[dtos.UserDto]{}).
		SetAuthToken(token).
		Get(getUsersEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*pagination.PageResultDto[dtos.UserDto])
	suite.NotNil(result)
	suite.Equal(1, result.TotalCount)
	suite.Equal(1, len(result.Items))
	suite.Equal(userConsts.DefaultAdminUserName, result.Items[0].UserName)
}

func (suite *UserTestSuite) TestShouldGetPaginatedUsers() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	userManager := userService.NewUserManager(suite.Pool)

	totalCount, err := userManager.GetUsersCount(suite.Ctx)
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetQueryParams(map[string]string{
			"sorting":        "user_name DESC",
			"skipCount":      "0",
			"maxResultCount": "1",
		}).
		SetResult(&pagination.PageResultDto[dtos.UserDto]{}).
		SetAuthToken(token).
		Get(getUsersEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*pagination.PageResultDto[dtos.UserDto])
	suite.NotNil(result)
	suite.Equal(totalCount, result.TotalCount)
	suite.Equal(1, len(result.Items))
	suite.Equal(userConsts.DefaultUserUserName, result.Items[0].UserName)
}

func (suite *UserTestSuite) TestShouldGetUsersUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetResult(&pagination.PageResultDto[dtos.UserDto]{}).
		SetAuthToken(token).
		Get(getUsersEndpoint)

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}
