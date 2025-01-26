package users

import (
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/danielgtaylor/huma/v2"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/internal/users/dtos"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

const (
	updateUserEndpoint = "/api/v1/user"
)

func (suite *UserTestSuite) TestShouldUpdateUser() {
	suite.ResetUsers()

	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	userManager := userService.NewUserManager(suite.Pool)
	newUser := GetFakeUser()
	err = userManager.CreateUser(suite.Ctx, newUser, "123123")
	suite.NoError(err)

	totalCount, err := userManager.GetUsersCount(suite.Ctx)
	suite.NoError(err)

	request := GetFakeEditUserDto()
	request.Id = &newUser.Id

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&GetUserDtoResult{}).
		SetAuthToken(token).
		Put(updateUserEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*GetUserDtoResult)
	suite.NotNil(result)
	suite.Equal(result.User.FirstName, request.FirstName)
	suite.Equal(result.User.LastName, request.LastName)
	suite.Equal(result.User.UserName, request.UserName)
	suite.Equal(result.User.Email, request.Email)

	updatedUser, err := userManager.GetUserById(suite.Ctx, result.User.Id)
	suite.NoError(err)
	suite.Equal(updatedUser.FirstName, request.FirstName)
	suite.Equal(updatedUser.LastName, request.LastName)
	suite.Equal(updatedUser.UserName, request.UserName)
	suite.Equal(updatedUser.NormalizedUserName, strings.ToUpper(request.UserName))
	suite.Equal(updatedUser.Email, request.Email)
	suite.Equal(updatedUser.NormalizedEmail, strings.ToUpper(request.Email))

	newTotalCount, err := userManager.GetUsersCount(suite.Ctx)
	suite.NoError(err)
	suite.Equal(totalCount, newTotalCount)
}

func (suite *UserTestSuite) TestShouldUpdateUserWithRoles() {
	suite.ResetUsers()

	adminToken, err := suite.LoginAsAdmin()
	suite.NoError(err)

	userManager := userService.NewUserManager(suite.Pool)
	newUser := GetFakeUser()
	err = userManager.CreateUser(suite.Ctx, newUser, "123123")
	suite.NoError(err)

	err = userManager.CreateUserRole(suite.Ctx, newUser.Id, 2)
	suite.NoError(err)

	userToken, err := suite.LoginAs(newUser.UserName)
	suite.NoError(err)

	request := GetFakeEditUserDto()
	request.Id = &newUser.Id

	request.RoleIds = make([]int64, 0)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&GetUserDtoResult{}).
		SetAuthToken(userToken).
		Put(updateUserEndpoint)

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())

	_, err = suite.CacheManager.Get(suite.Ctx, permissions.GenerateUserRoleCacheKey(newUser.Id))
	suite.NoError(err)

	resp, err = suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&GetUserDtoResult{}).
		SetAuthToken(adminToken).
		Put(updateUserEndpoint)

	suite.NoError(err)

	result := resp.Result().(*GetUserDtoResult)
	suite.NotNil(result)
	suite.Equal(result.User.FirstName, request.FirstName)
	suite.Equal(result.User.LastName, request.LastName)
	suite.Equal(result.User.UserName, request.UserName)
	suite.Equal(result.User.Email, request.Email)

	updatedUser, err := userManager.GetUserById(suite.Ctx, result.User.Id)
	suite.NoError(err)
	suite.Equal(updatedUser.FirstName, request.FirstName)
	suite.Equal(updatedUser.LastName, request.LastName)
	suite.Equal(updatedUser.UserName, request.UserName)
	suite.Equal(updatedUser.NormalizedUserName, strings.ToUpper(request.UserName))
	suite.Equal(updatedUser.Email, request.Email)
	suite.Equal(updatedUser.NormalizedEmail, strings.ToUpper(request.Email))

	_, err = suite.CacheManager.Get(suite.Ctx, permissions.GenerateUserRoleCacheKey(newUser.Id))
	suite.True(caching.CheckIsCacheValueNotFound(err))
}

func (suite *UserTestSuite) TestShouldUpdateUserWithError() {
	suite.ResetUsers()

	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		request          *dtos.EditUserDto
		wantErrorMessage string
		wantErrorCode    int
	}{
		{
			request:          GetUpdateUserWithErrorTestData(0),
			wantErrorMessage: "username admin is taken",
		},
		{
			request:          GetUpdateUserWithErrorTestData(1),
			wantErrorMessage: "email gktan@testgk.com is taken",
		},
		{
			request:          GetUpdateUserWithErrorTestData(2),
			wantErrorMessage: "Please enter the email",
		},
		{
			request:          GetUpdateUserWithErrorTestData(3),
			wantErrorMessage: "Please enter a valid email",
		},
		{
			request:          GetUpdateUserWithErrorTestData(4),
			wantErrorMessage: "Please enter the username",
		},
		{
			request:          GetUpdateUserWithErrorTestData(5),
			wantErrorMessage: "Invalid user id",
		},
		{
			request:          GetUpdateUserWithErrorTestData(6),
			wantErrorMessage: "Invalid user id",
		},
		{
			request:          GetUpdateUserWithErrorTestData(7),
			wantErrorMessage: "Invalid role id",
		},
		{
			request:          GetUpdateUserWithErrorTestData(8),
			wantErrorMessage: "user not found",
			wantErrorCode:    404,
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(tt.request).
			SetResult(&GetUserDtoResult{}).
			SetError(&huma.ErrorModel{}).
			SetAuthToken(token).
			Put(updateUserEndpoint)

		suite.NoError(err)
		if tt.wantErrorCode == 0 {
			suite.Equal(400, resp.StatusCode())
		} else {
			suite.Equal(tt.wantErrorCode, resp.StatusCode())
		}

		result := resp.Error().(*huma.ErrorModel)
		suite.NotNil(result)

		suite.Equal(tt.wantErrorMessage, result.Detail)
	}
}

func (suite *UserTestSuite) TestShouldUpdateUserUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	request := GetFakeEditUserDto()

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetAuthToken(token).
		Put(updateUserEndpoint)

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}

func GetUpdateUserWithErrorTestData(scenario int) *dtos.EditUserDto {
	userDto := GetFakeEditUserDto()
	userId := int64(2)
	userDto.Id = &userId

	switch scenario {
	case 0:
		userDto.UserName = userConsts.DefaultAdminUserName
	case 1:
		userId := int64(1)
		userDto.Id = &userId
		userDto.Email = "gktan@testgk.com"
	case 2:
		userDto.Email = ""
	case 3:
		userDto.Email = "invalid email"
	case 4:
		userDto.UserName = ""
	case 5:
		userId := int64(0)
		userDto.Id = &userId
	case 6:
		userDto.Id = nil
	case 7:
		userDto.RoleIds = make([]int64, 0)
		roleId := int64(20)
		userDto.RoleIds = append(userDto.RoleIds, roleId)
	case 8:
		userId := int64(100)
		userDto.Id = &userId
	}

	return userDto
}

func GetFakeEditUserDto() *dtos.EditUserDto {
	userId := int64(0)
	createUserDto := dtos.EditUserDto{}
	createUserDto.Id = &userId
	createUserDto.FirstName = gofakeit.FirstName()
	createUserDto.LastName = gofakeit.LastName()
	createUserDto.UserName = gofakeit.Username()
	createUserDto.Email = gofakeit.Email()
	createUserDto.Password = gofakeit.Password(true, true, true, true, false, 6)
	createUserDto.RoleIds = make([]int64, 0)

	return &createUserDto
}
