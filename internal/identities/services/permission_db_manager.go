package services

import (
	"context"

	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

type userPermissionDbManager struct {
	userRolePermissionManager userService.IUserRolePermissionManager
}

func NewPermissionDbManager(db userService.IUserRolePermissionManager) permissions.IPermissionDbManager {
	return &userPermissionDbManager{
		userRolePermissionManager: db,
	}
}

func (m *userPermissionDbManager) GetGrantedPermissionsFromDb(ctx context.Context, userId int64) (map[string]struct{}, error) {
	return m.userRolePermissionManager.SetUserPermissions(ctx, userId)
}
