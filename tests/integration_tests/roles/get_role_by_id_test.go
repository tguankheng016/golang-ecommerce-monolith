package roles

import (
	"fmt"

	roleConsts "github.com/tguankheng016/commerce-mono/internal/roles/constants"
	"github.com/tguankheng016/commerce-mono/internal/roles/dtos"
)

const (
	getRoleByIdEndpoint = "/api/v1/role"
)

type GetRoleByIdResult struct {
	Role dtos.CreateOrEditRoleDto
}

func (suite *RoleTestSuite) TestShouldGetCorrectRoleById() {
	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	tests := []struct {
		id           int64
		wantRolename string
	}{
		{
			id:           1,
			wantRolename: roleConsts.DefaultAdminRoleName,
		},
		{
			id:           0,
			wantRolename: "",
		},
	}

	for _, tt := range tests {
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetResult(&GetRoleByIdResult{}).
			SetAuthToken(token).
			Get(fmt.Sprintf("%s/%d", getRoleByIdEndpoint, tt.id))

		suite.NoError(err)
		suite.Equal(200, resp.StatusCode())

		result := resp.Result().(*GetRoleByIdResult)
		suite.NotNil(result)
		suite.NotNil(result.Role)
		suite.Equal(result.Role.Name, tt.wantRolename)

		if tt.id == 0 {
			suite.Nil(result.Role.Id)
		} else {
			suite.Equal(int64(1), *result.Role.Id)
		}
	}
}

func (suite *RoleTestSuite) TestShouldGetErrorWhenGetRoleById() {
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
			SetResult(&GetRoleByIdResult{}).
			SetAuthToken(token).
			Get(fmt.Sprintf("%s/%d", getRoleByIdEndpoint, tt.id))

		suite.NoError(err)
		suite.Equal(tt.wantErrorCode, resp.StatusCode())
	}
}

func (suite *RoleTestSuite) TestShouldGetRoleByIdUnauthorized() {
	token, err := suite.LoginAsUser()
	suite.NoError(err)

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetResult(&GetRoleByIdResult{}).
		SetAuthToken(token).
		Get(fmt.Sprintf("%s/%d", getRoleByIdEndpoint, 1))

	suite.NoError(err)
	suite.Equal(403, resp.StatusCode())
}
