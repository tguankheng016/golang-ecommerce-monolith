package permissions

import "fmt"

type UserPermissionCacheItem struct {
	UserId                int64
	Permissions           map[string]string
	ProhibitedPermissions map[string]string
}

const UserPermissionCacheName = "UserPermissions"

func NewUserPermissionCacheItem(userId int64, permissions map[string]string, prohibitedPermissions map[string]string) *UserPermissionCacheItem {
	return &UserPermissionCacheItem{
		UserId:                userId,
		Permissions:           permissions,
		ProhibitedPermissions: prohibitedPermissions,
	}
}

func GenerateUserPermissionCacheKey(userId int64) string {
	return fmt.Sprintf("%s:u%d", UserPermissionCacheName, userId)
}
