package users

import (
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/danielgtaylor/huma/v2"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/internal/users/dtos"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
)

const (
	createUserEndpoint = "/api/v1/user"
)

type GetUserDtoResult struct {
	User dtos.UserDto
}

func (suite *UserTestSuite) TestShouldCreateUser() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	userManager := userService.NewUserManager(suite.Pool)

	totalCount, err := userManager.GetUsersCount(suite.Ctx)
	suite.NoError(err)

	request := GetFakeCreateUserDto()

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&GetUserDtoResult{}).
		SetAuthToken(token).
		Post(createUserEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*GetUserDtoResult)
	suite.NotNil(result)
	suite.Equal(result.User.FirstName, request.FirstName)
	suite.Equal(result.User.LastName, request.LastName)
	suite.Equal(result.User.UserName, request.UserName)
	suite.Equal(result.User.Email, request.Email)

	newCreatedUser, err := userManager.GetUserById(suite.Ctx, result.User.Id)
	suite.NoError(err)
	suite.Equal(newCreatedUser.FirstName, request.FirstName)
	suite.Equal(newCreatedUser.LastName, request.LastName)
	suite.Equal(newCreatedUser.UserName, request.UserName)
	suite.Equal(newCreatedUser.NormalizedUserName, strings.ToUpper(request.UserName))
	suite.Equal(newCreatedUser.Email, request.Email)
	suite.Equal(newCreatedUser.NormalizedEmail, strings.ToUpper(request.Email))

	newTotalCount, err := userManager.GetUsersCount(suite.Ctx)
	suite.NoError(err)
	suite.Equal(totalCount+1, newTotalCount)
}

func (suite *UserTestSuite) TestShouldCreateUserWithError() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		request          *dtos.CreateUserDto
		wantErrorMessage string
	}{
		{
			request:          GetCreateUserWithErrorTestData(0),
			wantErrorMessage: "username admin is taken",
		},
		{
			request:          GetCreateUserWithErrorTestData(1),
			wantErrorMessage: "email gktan@testgk.com is taken",
		},
		{
			request:          GetCreateUserWithErrorTestData(2),
			wantErrorMessage: "Please enter the email",
		},
		{
			request:          GetCreateUserWithErrorTestData(3),
			wantErrorMessage: "Please enter a valid email",
		},
		{
			request:          GetCreateUserWithErrorTestData(4),
			wantErrorMessage: "Please enter the username",
		},
		{
			request:          GetCreateUserWithErrorTestData(5),
			wantErrorMessage: "Please enter the password",
		},
		{
			request:          GetCreateUserWithErrorTestData(6),
			wantErrorMessage: "Invalid user id",
		},
		{
			request:          GetCreateUserWithErrorTestData(7),
			wantErrorMessage: "Invalid user id",
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(tt.request).
			SetResult(&GetUserDtoResult{}).
			SetError(&huma.ErrorModel{}).
			SetAuthToken(token).
			Post(createUserEndpoint)

		suite.NoError(err)
		suite.Equal(400, resp.StatusCode())

		result := resp.Error().(*huma.ErrorModel)
		suite.NotNil(result)

		suite.Equal(tt.wantErrorMessage, result.Detail)
	}
}

func (suite *UserTestSuite) TestShouldCreateUserUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	request := GetFakeCreateUserDto()

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetAuthToken(token).
		Post(createUserEndpoint)

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}

func GetCreateUserWithErrorTestData(scenario int) *dtos.CreateUserDto {
	userDto := GetFakeCreateUserDto()
	switch scenario {
	case 0:
		userDto.UserName = userConsts.DefaultAdminUserName
	case 1:
		userDto.Email = "gktan@testgk.com"
	case 2:
		userDto.Email = ""
	case 3:
		userDto.Email = "invalid email"
	case 4:
		userDto.UserName = ""
	case 5:
		userDto.Password = ""
	case 6:
		userId := int64(1)
		userDto.Id = &userId
	case 7:
		userId := int64(-1)
		userDto.Id = &userId
	}

	return userDto
}

func GetFakeCreateUserDto() *dtos.CreateUserDto {
	userId := int64(0)
	createUserDto := dtos.CreateUserDto{}
	createUserDto.Id = &userId
	createUserDto.FirstName = gofakeit.FirstName()
	createUserDto.LastName = gofakeit.LastName()
	createUserDto.UserName = gofakeit.Username()
	createUserDto.Email = gofakeit.Email()
	createUserDto.Password = gofakeit.Password(true, true, true, true, false, 6)
	createUserDto.RoleIds = make([]int64, 0)

	return &createUserDto
}
