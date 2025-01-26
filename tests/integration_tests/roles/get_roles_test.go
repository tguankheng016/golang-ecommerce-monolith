package roles

import (
	roleConsts "github.com/tguankheng016/commerce-mono/internal/roles/constants"
	"github.com/tguankheng016/commerce-mono/internal/roles/dtos"
	roleService "github.com/tguankheng016/commerce-mono/internal/roles/services"
	"github.com/tguankheng016/commerce-mono/pkg/core/pagination"
)

const (
	getRolesEndpoint = "/api/v1/roles"
)

func (suite *RoleTestSuite) TestShouldGetRoles() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	roleManager := roleService.NewRoleManager(suite.Pool)

	totalCount, err := roleManager.GetRolesCount(suite.Ctx)
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetResult(&pagination.PageResultDto[dtos.RoleDto]{}).
		SetAuthToken(token).
		Get(getRolesEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*pagination.PageResultDto[dtos.RoleDto])
	suite.NotNil(result)
	suite.Equal(totalCount, result.TotalCount)
	suite.Equal(totalCount, len(result.Items))
}

func (suite *RoleTestSuite) TestShouldGetFilteredRoles() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetQueryParams(map[string]string{
			"filters": "admi",
		}).
		SetResult(&pagination.PageResultDto[dtos.RoleDto]{}).
		SetAuthToken(token).
		Get(getRolesEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*pagination.PageResultDto[dtos.RoleDto])
	suite.NotNil(result)
	suite.Equal(1, result.TotalCount)
	suite.Equal(1, len(result.Items))
	suite.Equal(roleConsts.DefaultAdminRoleName, result.Items[0].Name)
}

func (suite *RoleTestSuite) TestShouldGetPaginatedRoles() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	roleManager := roleService.NewRoleManager(suite.Pool)

	totalCount, err := roleManager.GetRolesCount(suite.Ctx)
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetQueryParams(map[string]string{
			"sorting":        "name DESC",
			"skipCount":      "0",
			"maxResultCount": "1",
		}).
		SetResult(&pagination.PageResultDto[dtos.RoleDto]{}).
		SetAuthToken(token).
		Get(getRolesEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*pagination.PageResultDto[dtos.RoleDto])
	suite.NotNil(result)
	suite.Equal(totalCount, result.TotalCount)
	suite.Equal(1, len(result.Items))
	suite.Equal(roleConsts.DefaultUserRoleName, result.Items[0].Name)
}

func (suite *RoleTestSuite) TestShouldGetRolesUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetResult(&pagination.PageResultDto[dtos.RoleDto]{}).
		SetAuthToken(token).
		Get(getRolesEndpoint)

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}
