package identities

import (
	"github.com/go-resty/resty/v2"
	get_current_session "github.com/tguankheng016/commerce-mono/internal/identities/features/getting_current_session/v1"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

const (
	getCurrentSessionEndpoint = "/api/v2/identities/current-session"
)

func (suite *IdentityTestSuite) TestShouldGetCorrectCurrentSession() {
	tests := []struct {
		username string
	}{
		{
			username: userConsts.DefaultAdminUserName,
		},
		{
			username: userConsts.DefaultUserUserName,
		},
		{
			username: "",
		},
	}

	for _, tt := range tests {
		var resp *resty.Response
		var err error

		if tt.username != "" {
			token, err := suite.LoginAs(tt.username)
			suite.NoError(err)

			resp, err = suite.Client.R().
				SetContext(suite.Ctx).
				SetResult(&get_current_session.GetCurrentSessionResult{}).
				SetAuthToken(token).
				Get(getCurrentSessionEndpoint)

			suite.NoError(err)
		} else {
			resp, err = suite.Client.R().
				SetContext(suite.Ctx).
				SetResult(&get_current_session.GetCurrentSessionResult{}).
				Get(getCurrentSessionEndpoint)

			suite.NoError(err)
		}

		suite.Equal(200, resp.StatusCode())

		result := resp.Result().(*get_current_session.GetCurrentSessionResult)
		suite.NotNil(result)

		allPermissions := permissions.GetAppPermissions().Items
		suite.Equal(len(allPermissions), len(result.AllPermissions))

		if tt.username == "" {
			suite.Nil(result.User)
		} else {
			suite.Equal(tt.username, result.User.UserName)

			if tt.username == userConsts.DefaultAdminUserName {
				suite.Equal(len(allPermissions), len(result.AllPermissions))
			} else {
				suite.Equal(0, len(result.GrantedPermissions))
			}
		}
	}
}
