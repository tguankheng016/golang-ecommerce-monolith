package roles

import (
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/danielgtaylor/huma/v2"
	roleConsts "github.com/tguankheng016/commerce-mono/internal/roles/constants"
	"github.com/tguankheng016/commerce-mono/internal/roles/dtos"
	roleService "github.com/tguankheng016/commerce-mono/internal/roles/services"
	appPermissions "github.com/tguankheng016/commerce-mono/pkg/permissions"
)

const (
	updateRoleEndpoint = "/api/v1/role"
)

func (suite *RoleTestSuite) TestShouldUpdateRole() {
	suite.ResetRoles()

	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	roleManager := roleService.NewRoleManager(suite.Pool)
	newRole := GetFakeRole()
	err = roleManager.CreateRole(suite.Ctx, newRole)
	suite.NoError(err)

	totalCount, err := roleManager.GetRolesCount(suite.Ctx)
	suite.NoError(err)

	request := GetFakeEditRoleDto()
	request.Id = &newRole.Id

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&GetRoleDtoResult{}).
		SetAuthToken(token).
		Put(updateRoleEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*GetRoleDtoResult)
	suite.NotNil(result)
	suite.Equal(result.Role.Name, request.Name)

	updatedRole, err := roleManager.GetRoleById(suite.Ctx, result.Role.Id)
	suite.NoError(err)
	suite.Equal(updatedRole.Name, request.Name)
	suite.Equal(updatedRole.NormalizedName, strings.ToUpper(request.Name))
	suite.Equal(updatedRole.IsDefault, request.IsDefault)

	newTotalCount, err := roleManager.GetRolesCount(suite.Ctx)
	suite.NoError(err)
	suite.Equal(totalCount, newTotalCount)
}

func (suite *RoleTestSuite) TestShouldUpdateRoleWithPermissions() {
	suite.ResetRoles()
	suite.ResetRolePermissions()

	roleManager := roleService.NewRoleManager(suite.Pool)

	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		roleId                       int64
		permissionToUpdate           []string
		expectedGrantedPermission    []string
		expectedProhibitedPermission []string
	}{
		{
			roleId:                       1,
			permissionToUpdate:           GetUpdateRoleWithPermissionsTestData(0),
			expectedGrantedPermission:    GetExpectedGrantedPermissionsForUpdateRoleWithPermissions(0),
			expectedProhibitedPermission: GetExpectedProhitbitedPermissionsForUpdateRoleWithPermissions(0),
		},
		{
			roleId:                       2,
			permissionToUpdate:           GetUpdateRoleWithPermissionsTestData(1),
			expectedGrantedPermission:    GetExpectedGrantedPermissionsForUpdateRoleWithPermissions(1),
			expectedProhibitedPermission: GetExpectedProhitbitedPermissionsForUpdateRoleWithPermissions(1),
		},
	}

	for _, tt := range tests {
		role, err := roleManager.GetRoleById(suite.Ctx, tt.roleId)
		suite.NoError(err)
		suite.NotNil(role)

		request := dtos.EditRoleDto{}
		request.Id = &role.Id
		request.Name = role.Name
		request.GrantedPermissions = tt.permissionToUpdate

		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(request).
			SetResult(&GetRoleDtoResult{}).
			SetError(&huma.ErrorModel{}).
			SetAuthToken(token).
			Put(updateRoleEndpoint)

		suite.NoError(err)
		suite.Equal(200, resp.StatusCode())

		result := resp.Result().(*GetRoleDtoResult)
		suite.NotNil(result)

		for _, grantedPermission := range tt.expectedGrantedPermission {
			permission, err := roleManager.GetRolePermission(suite.Ctx, role.Id, grantedPermission)
			suite.NoError(err)
			suite.NotNil(permission)
			suite.Equal(grantedPermission, permission.Name)
			suite.True(permission.IsGranted)
		}

		for _, prohibitedPermission := range tt.expectedProhibitedPermission {
			permission, err := roleManager.GetRolePermission(suite.Ctx, role.Id, prohibitedPermission)
			suite.NoError(err)
			suite.NotNil(permission)
			suite.Equal(prohibitedPermission, permission.Name)
			suite.False(permission.IsGranted)
		}
	}

	suite.ResetRolePermissions()
}

func (suite *RoleTestSuite) TestShouldUpdateRoleWithError() {
	suite.ResetRoles()
	suite.ResetRolePermissions()

	roleManager := roleService.NewRoleManager(suite.Pool)
	newRole := GetFakeRole()
	err := roleManager.CreateRole(suite.Ctx, newRole)
	suite.NoError(err)

	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		request          *dtos.EditRoleDto
		wantErrorMessage string
		wantErrorCode    int
	}{
		{
			request:          GetUpdateRoleWithErrorTestData(0, newRole.Id),
			wantErrorMessage: "role name Admin is taken",
		},
		{
			request:          GetUpdateRoleWithErrorTestData(1, newRole.Id),
			wantErrorMessage: "invalid permission name",
		},
		{
			request:          GetUpdateRoleWithErrorTestData(2, newRole.Id),
			wantErrorMessage: "Invalid role id",
		},
		{
			request:          GetUpdateRoleWithErrorTestData(3, newRole.Id),
			wantErrorMessage: "Invalid role id",
		},
		{
			request:          GetUpdateRoleWithErrorTestData(4, newRole.Id),
			wantErrorMessage: "role not found",
			wantErrorCode:    404,
		},
		{
			request:          GetUpdateRoleWithErrorTestData(5, newRole.Id),
			wantErrorMessage: "Please enter the role name",
		},
		{
			request:          GetUpdateRoleWithErrorTestData(6, newRole.Id),
			wantErrorMessage: "You cannot change the name of static role",
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(tt.request).
			SetResult(&GetRoleDtoResult{}).
			SetError(&huma.ErrorModel{}).
			SetAuthToken(token).
			Put(updateRoleEndpoint)

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

func (suite *RoleTestSuite) TestShouldUpdateRoleUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	request := GetFakeEditRoleDto()

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetAuthToken(token).
		Put(updateRoleEndpoint)

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}

func GetUpdateRoleWithErrorTestData(scenario int, newRoleId int64) *dtos.EditRoleDto {
	roleDto := GetFakeEditRoleDto()
	roleDto.Id = &newRoleId

	switch scenario {
	case 0:
		roleDto.Name = roleConsts.DefaultAdminRoleName
	case 1:
		roleDto.GrantedPermissions = make([]string, 0)
		roleDto.GrantedPermissions = append(roleDto.GrantedPermissions, "Test")
	case 2:
		roleId := int64(0)
		roleDto.Id = &roleId
	case 3:
		roleDto.Id = nil
	case 4:
		roleId := int64(100)
		roleDto.Id = &roleId
	case 5:
		roleDto.Name = ""
	case 6:
		roleId := int64(2)
		roleDto.Id = &roleId
		roleDto.Name = "Test222"
	}

	return roleDto
}

func GetUpdateRoleWithPermissionsTestData(scenario int) []string {
	permissions := make([]string, 0)

	switch scenario {
	case 0:
		// Admin With Prohibited
		allPermissions := appPermissions.GetAppPermissions().Items
		for _, allPermission := range allPermissions {
			if allPermission.Name != appPermissions.PagesAdministrationUsers {
				permissions = append(permissions, allPermission.Name)
			}
		}
	case 1:
		/// User Without Prohibited
		permissions = append(permissions, appPermissions.PagesAdministrationUsers)
		permissions = append(permissions, appPermissions.PagesAdministrationUsersCreate)
	}

	return permissions
}

func GetExpectedGrantedPermissionsForUpdateRoleWithPermissions(scenario int) []string {
	permissions := make([]string, 0)

	switch scenario {
	case 0:
		// Admin With Prohibited
		// return empty because admin is granted by default
	case 1:
		/// User Without Prohibited
		permissions = append(permissions, appPermissions.PagesAdministrationUsers)
		permissions = append(permissions, appPermissions.PagesAdministrationUsersCreate)
	}

	return permissions
}

func GetExpectedProhitbitedPermissionsForUpdateRoleWithPermissions(scenario int) []string {
	permissions := make([]string, 0)

	switch scenario {
	case 0:
		// Admin With Prohibited
		permissions = append(permissions, appPermissions.PagesAdministrationUsers)
	case 1:
		/// User Without Prohibited
	}

	return permissions
}

func GetFakeEditRoleDto() *dtos.EditRoleDto {
	roleId := int64(0)
	createRoleDto := dtos.EditRoleDto{}
	createRoleDto.Id = &roleId
	createRoleDto.Name = gofakeit.BeerName()

	return &createRoleDto
}
