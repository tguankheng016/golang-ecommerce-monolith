package permissions

import "fmt"

type UserRoleCacheItem struct {
	UserId  int64
	RoleIds []int64
}

const UserRoleCacheName = "UserRoles"

func NewUserRoleCacheItem(userId int64, roleIds []int64) *UserRoleCacheItem {
	return &UserRoleCacheItem{
		UserId:  userId,
		RoleIds: roleIds,
	}
}

func GenerateUserRoleCacheKey(userId int64) string {
	return fmt.Sprintf("%s:u%d", UserRoleCacheName, userId)
}
