package roles

import (
	"fmt"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	roleService "github.com/tguankheng016/commerce-mono/internal/roles/services"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

const (
	deleteRoleEndpoint = "/api/v1/role"
)

func (suite *RoleTestSuite) TestShouldDeleteRole() {
	suite.ResetRoles()
	suite.ResetRolePermissions()

	roleManager := roleService.NewRoleManager(suite.Pool)
	newRole := GetFakeRole()
	err := roleManager.CreateRole(suite.Ctx, newRole)
	suite.NoError(err)

	err = roleManager.CreateRolePermission(suite.Ctx, newRole.Id, permissions.PagesAdministrationUsers, true)
	suite.NoError(err)

	userManager := userService.NewUserManager(suite.Pool)
	err = userManager.CreateUserRole(suite.Ctx, 2, newRole.Id)
	suite.NoError(err)

	// Login As New Role To Trigger Middleware
	userToken, err := suite.LoginAsUser()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetAuthToken(userToken).
		Delete(fmt.Sprintf("%s/%d", deleteRoleEndpoint, newRole.Id))

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())

	// Ensure User Role Is Cached
	userRolesCache, err := suite.CacheManager.Get(suite.Ctx, permissions.GenerateUserRoleCacheKey(2))
	suite.NoError(err)
	suite.NotEmpty(userRolesCache)

	// Login And Perform Delete As Admin
	adminToken, err := suite.LoginAsAdmin()
	suite.NoError(err)

	resp, err = suite.Client.R().
		SetContext(suite.Ctx).
		SetAuthToken(adminToken).
		Delete(fmt.Sprintf("%s/%d", deleteRoleEndpoint, newRole.Id))

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	// Ensure role is deleted
	role, err := roleManager.GetRoleById(suite.Ctx, newRole.Id)
	suite.NoError(err)
	suite.Nil(role)

	// Ensure user role is removed
	roleIds, err := userManager.GetUserRoleIds(suite.Ctx, 2)
	suite.NoError(err)
	suite.True(!slices.Contains(roleIds, newRole.Id))

	// Ensure role permission is removed
	rolePermission, err := roleManager.GetRolePermission(suite.Ctx, newRole.Id, permissions.PagesAdministrationUsers)
	suite.NoError(err)
	suite.Nil(rolePermission)

	// Ensure user role cached is clear
	_, err = suite.CacheManager.Get(suite.Ctx, permissions.GenerateUserRoleCacheKey(2))
	suite.True(caching.CheckIsCacheValueNotFound(err))
}

func (suite *RoleTestSuite) TestShouldDeleteRoleWithError() {
	suite.ResetRoles()
	suite.ResetRolePermissions()

	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		deleteRoleId     int64
		wantErrorMessage string
		wantErrorCode    int
	}{
		{
			deleteRoleId:     1,
			wantErrorMessage: "You cannot delete static role!",
		},
		{
			deleteRoleId:     100,
			wantErrorMessage: "role not found",
			wantErrorCode:    404,
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetError(&huma.ErrorModel{}).
			SetAuthToken(token).
			Delete(fmt.Sprintf("%s/%d", deleteRoleEndpoint, tt.deleteRoleId))

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

func (suite *RoleTestSuite) TestShouldDeleteRoleUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetError(&huma.ErrorModel{}).
		SetAuthToken(token).
		Delete(fmt.Sprintf("%s/%d", deleteRoleEndpoint, 1))

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}
