package users

import (
	"fmt"
	"slices"

	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	appPermissions "github.com/tguankheng016/commerce-mono/pkg/permissions"
)

type GetUserPermissionsResult struct {
	Items []string
}

type UpdateUserPermissionRequest struct {
	GrantedPermissions []string
}

func (suite *UserTestSuite) TestUserWithPermissions() {
	suite.ResetUsers()
	suite.ResetUserPermissions()

	token, err := suite.LoginAsAdmin()
	suite.NoError(err)

	userManager := userService.NewUserManager(suite.Pool)

	tests := []struct {
		userId                        int64
		updatedPermissions            []string
		expectedGetPermissions        []string
		expectedGrantedPermissions    []string
		expectedProhibitedPermissions []string
	}{
		{
			userId:                        1,
			updatedPermissions:            GetUpdatedUserPermissions(1),
			expectedGetPermissions:        GetExpectedGetUserPermissions(1),
			expectedGrantedPermissions:    GetExpectedGrantedUserPermissions(1),
			expectedProhibitedPermissions: GetExpectedProhibitedUserPermissions(1),
		},
	}

	for _, tt := range tests {
		// Get User Permissions
		resp, err := suite.Client.R().
			SetContext(suite.Ctx).
			SetResult(&GetUserPermissionsResult{}).
			SetAuthToken(token).
			Get(fmt.Sprintf("/api/v1/user/%d/permissions", tt.userId))

		suite.NoError(err)
		suite.Equal(200, resp.StatusCode())

		result := resp.Result().(*GetUserPermissionsResult)
		suite.NotNil(result)

		suite.Equal(len(result.Items), len(tt.expectedGetPermissions))

		for _, expectedPermission := range tt.expectedGetPermissions {
			suite.True(slices.Contains(result.Items, expectedPermission))
		}

		// Update User Permissions
		updatedUserPermissionRequest := UpdateUserPermissionRequest{}
		updatedUserPermissionRequest.GrantedPermissions = tt.updatedPermissions
		resp, err = suite.Client.R().
			SetContext(suite.Ctx).
			SetBody(updatedUserPermissionRequest).
			SetAuthToken(token).
			Put(fmt.Sprintf("/api/v1/user/%d/permissions", tt.userId))

		suite.NoError(err)
		suite.Equal(200, resp.StatusCode())

		isGranted := true
		dbGrantedPermissions, err := userManager.GetUserPermissions(suite.Ctx, tt.userId, &isGranted)

		suite.NoError(err)
		suite.Equal(len(dbGrantedPermissions), len(tt.expectedGrantedPermissions))

		for _, dbGrantedPermission := range dbGrantedPermissions {
			suite.True(slices.Contains(tt.expectedGrantedPermissions, dbGrantedPermission.Name))
		}

		isGranted = false
		dbProhibitedPermissions, err := userManager.GetUserPermissions(suite.Ctx, tt.userId, &isGranted)

		suite.NoError(err)
		suite.Equal(len(dbProhibitedPermissions), len(tt.expectedProhibitedPermissions))

		for _, dbProhibitedPermission := range dbProhibitedPermissions {
			fmt.Print(slices.Contains(tt.expectedProhibitedPermissions, dbProhibitedPermission.Name))
			suite.True(slices.Contains(tt.expectedProhibitedPermissions, dbProhibitedPermission.Name))
		}

		// Reset User Permissions
		resp, err = suite.Client.R().
			SetContext(suite.Ctx).
			SetAuthToken(token).
			Put(fmt.Sprintf("/api/v1/user/%d/reset-permissions", tt.userId))

		suite.NoError(err)
		suite.Equal(200, resp.StatusCode())

		dbUserPermissions, err := userManager.GetUserPermissions(suite.Ctx, tt.userId, nil)

		suite.NoError(err)
		suite.Equal(0, len(dbUserPermissions))
	}
}

func GetExpectedGetUserPermissions(userId int64) []string {
	permissions := make([]string, 0)

	switch userId {
	case 1:
		// Admin
		allPermissions := appPermissions.GetAppPermissions().Items
		for _, allPermission := range allPermissions {
			permissions = append(permissions, allPermission.Name)
		}
	case 2:
		/// User
	}

	return permissions
}

func GetUpdatedUserPermissions(userId int64) []string {
	permissions := make([]string, 0)

	switch userId {
	case 1:
		// Admin
		allPermissions := appPermissions.GetAppPermissions().Items
		for _, allPermission := range allPermissions {
			if allPermission.Name != appPermissions.PagesAdministrationRolesDelete {
				permissions = append(permissions, allPermission.Name)
			}
		}
	case 2:
		/// User
		permissions = append(permissions, appPermissions.PagesAdministrationUsers)
		permissions = append(permissions, appPermissions.PagesAdministrationRoles)
	}

	return permissions
}

func GetExpectedGrantedUserPermissions(userId int64) []string {
	permissions := make([]string, 0)

	switch userId {
	case 1:
		// Admin
	case 2:
		/// User
		return GetUpdatedUserPermissions(userId)
	}

	return permissions
}

func GetExpectedProhibitedUserPermissions(userId int64) []string {
	permissions := make([]string, 0)

	switch userId {
	case 1:
		// Admin
		permissions = append(permissions, appPermissions.PagesAdministrationRolesDelete)
	case 2:
		/// User
	}

	return permissions
}
