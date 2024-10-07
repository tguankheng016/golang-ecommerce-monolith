package permissions

import "fmt"

type RolePermissionCacheItem struct {
	RoleId      int64
	Permissions map[string]string
}

const RolePermissionCacheName = "RolePermissions"

func NewRolePermissionCacheItem(roleId int64, permissions map[string]string) *RolePermissionCacheItem {
	return &RolePermissionCacheItem{
		RoleId:      roleId,
		Permissions: permissions,
	}
}

func GenerateRolePermissionCacheKey(roleId int64) string {
	return fmt.Sprintf("%s:r%d", RolePermissionCacheName, roleId)
}
