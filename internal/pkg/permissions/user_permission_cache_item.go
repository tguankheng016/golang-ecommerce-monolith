package permissions

import "fmt"

type UserPermissionCacheItem struct {
	UserId                int64
	Permissions           map[string]struct{}
	ProhibitedPermissions map[string]struct{}
}

const UserPermissionCacheName = "UserPermissions"

func NewUserPermissionCacheItem(userId int64, permissions map[string]struct{}, prohibitedPermissions map[string]struct{}) *UserPermissionCacheItem {
	return &UserPermissionCacheItem{
		UserId:                userId,
		Permissions:           permissions,
		ProhibitedPermissions: prohibitedPermissions,
	}
}

func GenerateUserPermissionCacheKey(userId int64) string {
	return fmt.Sprintf("%s:u%d", UserPermissionCacheName, userId)
}
