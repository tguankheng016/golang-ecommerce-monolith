package roles

import (
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/danielgtaylor/huma/v2"
	roleConsts "github.com/tguankheng016/commerce-mono/internal/roles/constants"
	"github.com/tguankheng016/commerce-mono/internal/roles/dtos"
	roleService "github.com/tguankheng016/commerce-mono/internal/roles/services"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

const (
	createRoleEndpoint = "/api/v1/role"
)

type GetRoleDtoResult struct {
	Role dtos.RoleDto
}

func (suite *RoleTestSuite) TestShouldCreateRole() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	roleManager := roleService.NewRoleManager(suite.Pool)

	totalCount, err := roleManager.GetRolesCount(suite.Ctx)
	suite.NoError(err)

	request := GetFakeCreateRoleDto()
	request.GrantedPermissions = append(request.GrantedPermissions, permissions.PagesAdministrationUsers)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&GetRoleDtoResult{}).
		SetAuthToken(token).
		Post(createRoleEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*GetRoleDtoResult)
	suite.NotNil(result)
	suite.Equal(result.Role.Name, request.Name)

	newCreatedRole, err := roleManager.GetRoleById(suite.Ctx, result.Role.Id)
	suite.NoError(err)
	suite.Equal(newCreatedRole.Name, request.Name)
	suite.Equal(newCreatedRole.NormalizedName, strings.ToUpper(request.Name))
	suite.Equal(newCreatedRole.IsDefault, request.IsDefault)

	newTotalCount, err := roleManager.GetRolesCount(suite.Ctx)
	suite.NoError(err)
	suite.Equal(totalCount+1, newTotalCount)

	rolePermission, err := roleManager.GetRolePermission(suite.Ctx, result.Role.Id, permissions.PagesAdministrationUsers)
	suite.NoError(err)
	suite.NotNil(rolePermission)
	suite.True(rolePermission.RoleId.Valid)
	suite.Equal(result.Role.Id, rolePermission.RoleId.Int64)
	suite.Equal(permissions.PagesAdministrationUsers, rolePermission.Name)
}

func (suite *RoleTestSuite) TestShouldCreateRoleWithError() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		request          *dtos.CreateRoleDto
		wantErrorMessage string
	}{
		{
			request:          GetCreateRoleWithErrorTestData(0),
			wantErrorMessage: "role name Admin is taken",
		},
		{
			request:          GetCreateRoleWithErrorTestData(1),
			wantErrorMessage: "invalid permission name",
		},
		{
			request:          GetCreateRoleWithErrorTestData(2),
			wantErrorMessage: "Invalid role id",
		},
		{
			request:          GetCreateRoleWithErrorTestData(3),
			wantErrorMessage: "Invalid role id",
		},
		{
			request:          GetCreateRoleWithErrorTestData(4),
			wantErrorMessage: "Please enter the role name",
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(tt.request).
			SetResult(&GetRoleDtoResult{}).
			SetError(&huma.ErrorModel{}).
			SetAuthToken(token).
			Post(createRoleEndpoint)

		suite.NoError(err)
		suite.Equal(400, resp.StatusCode())

		result := resp.Error().(*huma.ErrorModel)
		suite.NotNil(result)

		suite.Equal(tt.wantErrorMessage, result.Detail)
	}
}

func (suite *RoleTestSuite) TestShouldCreateRoleUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	request := GetFakeCreateRoleDto()

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetAuthToken(token).
		Post(createRoleEndpoint)

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}

func GetCreateRoleWithErrorTestData(scenario int) *dtos.CreateRoleDto {
	roleDto := GetFakeCreateRoleDto()
	switch scenario {
	case 0:
		roleDto.Name = roleConsts.DefaultAdminRoleName
	case 1:
		roleDto.GrantedPermissions = make([]string, 0)
		roleDto.GrantedPermissions = append(roleDto.GrantedPermissions, "Test")
	case 2:
		roleId := int64(1)
		roleDto.Id = &roleId
	case 3:
		roleId := int64(-1)
		roleDto.Id = &roleId
	case 4:
		roleDto.Name = ""
	}

	return roleDto
}

func GetFakeCreateRoleDto() *dtos.CreateRoleDto {
	roleId := int64(0)
	createRoleDto := dtos.CreateRoleDto{}
	createRoleDto.Id = &roleId
	createRoleDto.Name = gofakeit.BeerName()

	return &createRoleDto
}
