package users

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
)

const (
	deleteUserEndpoint = "/api/v1/user"
)

func (suite *UserTestSuite) TestShouldDeleteUser() {
	suite.ResetUsers()

	userManager := userService.NewUserManager(suite.Pool)
	newUser := GetFakeUser()
	err := userManager.CreateUser(suite.Ctx, newUser, "123123")
	suite.NoError(err)

	err = userManager.CreateUserRole(suite.Ctx, newUser.Id, 2)
	suite.NoError(err)

	err = userManager.CreateUserPermission(suite.Ctx, newUser.Id, permissions.PagesAdministrationUsers, true)
	suite.NoError(err)

	// Login As New User To Trigger Middleware
	userToken, err := suite.LoginAs(newUser.UserName)
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetAuthToken(userToken).
		Delete(fmt.Sprintf("%s/%d", deleteUserEndpoint, newUser.Id))

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())

	// Ensure User Stamped Is Cached
	securityStampCache, err := suite.CacheManager.Get(suite.Ctx, jwt.GenerateStampCacheKey(newUser.Id))
	suite.NoError(err)
	suite.NotEmpty(securityStampCache)

	// Login And Perform Delete As Admin
	adminToken, err := suite.LoginAsAdmin()
	suite.NoError(err)

	resp, err = suite.Client.R().
		SetContext(suite.Ctx).
		SetAuthToken(adminToken).
		Delete(fmt.Sprintf("%s/%d", deleteUserEndpoint, newUser.Id))

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	// Ensure user is deleted
	user, err := userManager.GetUserById(suite.Ctx, newUser.Id)
	suite.NoError(err)
	suite.Nil(user)

	// Ensure user role is removed
	roleIds, err := userManager.GetUserRoleIds(suite.Ctx, newUser.Id)
	suite.NoError(err)
	suite.Equal(0, len(roleIds))

	// Ensure user permission is removed
	userPermission, err := userManager.GetUserPermission(suite.Ctx, newUser.Id, permissions.PagesAdministrationUsers)
	suite.NoError(err)
	suite.Nil(userPermission)

	// Ensure user cached is clear
	_, err = suite.CacheManager.Get(suite.Ctx, jwt.GenerateStampCacheKey(newUser.Id))
	suite.True(caching.CheckIsCacheValueNotFound(err))
}

func (suite *UserTestSuite) TestShouldDeleteUserWithError() {
	suite.ResetUsers()

	userManager := userService.NewUserManager(suite.Pool)
	newUser := GetFakeUser()
	err := userManager.CreateUser(suite.Ctx, newUser, "123123")
	suite.NoError(err)

	err = userManager.CreateUserRole(suite.Ctx, newUser.Id, 2)
	suite.NoError(err)

	err = userManager.CreateUserPermission(suite.Ctx, newUser.Id, permissions.PagesAdministrationUsersDelete, true)
	suite.NoError(err)

	token, err := suite.LoginAs(newUser.UserName)
	suite.NoError(err)

	tests := []struct {
		deleteUserId     int64
		wantErrorMessage string
		wantErrorCode    int
	}{
		{
			deleteUserId:     1,
			wantErrorMessage: "You cannot delete admin's account!",
		},
		{
			deleteUserId:     newUser.Id,
			wantErrorMessage: "You cannot delete your own account!",
		},
		{
			deleteUserId:     100,
			wantErrorMessage: "user not found",
			wantErrorCode:    404,
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetError(&huma.ErrorModel{}).
			SetAuthToken(token).
			Delete(fmt.Sprintf("%s/%d", deleteUserEndpoint, tt.deleteUserId))

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

func (suite *UserTestSuite) TestShouldDeleteUserUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetError(&huma.ErrorModel{}).
		SetAuthToken(token).
		Delete(fmt.Sprintf("%s/%d", deleteUserEndpoint, 1))

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}
