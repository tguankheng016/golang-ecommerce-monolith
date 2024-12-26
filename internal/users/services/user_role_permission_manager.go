package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	roleConsts "github.com/tguankheng016/commerce-mono/internal/roles/constants"
	roleService "github.com/tguankheng016/commerce-mono/internal/roles/services"
	"github.com/tguankheng016/commerce-mono/internal/users/models"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"go.uber.org/zap"
)

const (
	DefaultCacheExpiration = 1 * time.Hour
)

type IUserRolePermissionManager interface {
	IsGranted(ctx context.Context, userId int64, permissionName string) (bool, error)
	SetUserPermissions(ctx context.Context, userId int64) (map[string]struct{}, error)
	SetRolePermissions(ctx context.Context, roleId int64) (map[string]struct{}, error)
	RemoveUserRoleCaches(ctx context.Context, userId int64)
}

type userRolePermissionManager struct {
	db           postgres.IPgxDbConn
	cacheManager *cache.Cache[string]
}

func NewUserRolePermissionManager(db *pgxpool.Pool, cacheManager *cache.Cache[string]) IUserRolePermissionManager {
	return &userRolePermissionManager{
		db:           db,
		cacheManager: cacheManager,
	}
}

func (u *userRolePermissionManager) IsGranted(ctx context.Context, userId int64, permissionName string) (bool, error) {
	grantedPermissions, err := u.SetUserPermissions(ctx, userId)
	if err != nil {
		return false, err
	}

	_, ok := grantedPermissions[permissionName]
	return ok, nil
}

func (u *userRolePermissionManager) SetUserPermissions(ctx context.Context, userId int64) (map[string]struct{}, error) {
	userManager := NewUserManager(u.db)
	roleIds, err := userManager.GetUserRoleIds(ctx, userId)
	if err != nil {
		return nil, err
	}

	userRoleCacheItem, err := caching.MarshalCacheItem(permissions.NewUserRoleCacheItem(userId, roleIds))
	if err != nil {
		return nil, err
	}
	if err := u.cacheManager.Set(ctx, permissions.GenerateUserRoleCacheKey(userId), string(userRoleCacheItem), store.WithExpiration(DefaultCacheExpiration)); err != nil {
		// Dont return just log
		logging.Logger.Error("error in setting user role caches", zap.Error(err))
	}

	userPermissions, err := u.getUserPermissions(ctx, userId)
	if err != nil {
		return nil, err
	}

	grantedPermissions := make(map[string]struct{})
	grantedUserPermissions := make(map[string]struct{})
	prohibitedUserPermissions := make(map[string]struct{})
	for _, permission := range userPermissions {
		if permission.IsGranted {
			grantedUserPermissions[permission.Name] = struct{}{}
			grantedPermissions[permission.Name] = struct{}{}
		} else {
			prohibitedUserPermissions[permission.Name] = struct{}{}
		}
	}

	userPermissionCacheItem, err := caching.MarshalCacheItem(permissions.NewUserPermissionCacheItem(userId, grantedUserPermissions, prohibitedUserPermissions))
	if err != nil {
		return nil, err
	}
	if err := u.cacheManager.Set(ctx, permissions.GenerateUserPermissionCacheKey(userId), string(userPermissionCacheItem), store.WithExpiration(DefaultCacheExpiration)); err != nil {
		// Dont return just log
		logging.Logger.Error("error in setting user permissions caches", zap.Error(err))
	}

	for _, roleId := range roleIds {
		rolePermissions, err := u.SetRolePermissions(ctx, roleId)
		if err != nil {
			return nil, err
		}

		for permission := range rolePermissions {
			if _, ok := prohibitedUserPermissions[permission]; !ok {
				if _, ok := grantedPermissions[permission]; !ok {
					// key does not exists
					grantedPermissions[permission] = struct{}{}
				}
			}
		}
	}

	return grantedPermissions, nil
}

func (u *userRolePermissionManager) SetRolePermissions(ctx context.Context, roleId int64) (map[string]struct{}, error) {
	roleManager := roleService.NewRoleManager(u.db)
	role, err := roleManager.GetRoleById(ctx, roleId)
	if err != nil {
		return nil, err
	}

	isAdmin := strings.EqualFold(role.Name, roleConsts.DefaultAdminRoleName)

	rolePermissions, err := u.getRolePermissions(ctx, roleId)
	if err != nil {
		return nil, err
	}

	if isAdmin {
		allPermissions := permissions.GetAppPermissions().Items

		// Get all prohibited permissions for admin role
		prohibitedPermissions := make(map[string]struct{})
		for _, permission := range rolePermissions {
			if !permission.IsGranted {
				prohibitedPermissions[permission.Name] = struct{}{}
			}
		}

		// Excluded prohibited permissions
		grantedPermissions := make(map[string]struct{})
		for key := range allPermissions {
			if _, ok := prohibitedPermissions[key]; !ok {
				grantedPermissions[key] = struct{}{}
			}
		}

		rolePermissionCacheItem, err := caching.MarshalCacheItem(permissions.NewRolePermissionCacheItem(roleId, grantedPermissions))
		if err != nil {
			return nil, err
		}
		if err := u.cacheManager.Set(ctx, permissions.GenerateRolePermissionCacheKey(roleId), string(rolePermissionCacheItem), store.WithExpiration(DefaultCacheExpiration)); err != nil {
			// Dont return just log
			logging.Logger.Error("error in setting role permission caches", zap.Error(err))
		}

		return grantedPermissions, nil
	} else {
		grantedPermissions := make(map[string]struct{})
		for _, permission := range rolePermissions {
			if permission.IsGranted {
				grantedPermissions[permission.Name] = struct{}{}
			}
		}

		return grantedPermissions, nil
	}
}

func (u *userRolePermissionManager) RemoveUserRoleCaches(ctx context.Context, userId int64) {
	if err := u.cacheManager.Delete(ctx, permissions.GenerateUserRoleCacheKey(userId)); err != nil {
		logging.Logger.Error("error in removing user roles caches", zap.Error(err))
	}
}

func (u *userRolePermissionManager) getUserPermissions(ctx context.Context, userId int64) ([]models.UserRolePermission, error) {
	query := `
		select user_role_permissions.* 
		from users join user_role_permissions on users.id = user_role_permissions.user_id 
		where users.is_deleted = false and users.id = @userId
	`

	args := pgx.NamedArgs{
		"userId": userId,
	}

	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user role permissions: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.UserRolePermission])
}

func (u *userRolePermissionManager) getRolePermissions(ctx context.Context, roleId int64) ([]models.UserRolePermission, error) {
	query := `
		select user_role_permissions.* 
		from roles join user_role_permissions on roles.id = user_role_permissions.role_id 
		where roles.is_deleted = false and roles.id = @roleId
	`

	args := pgx.NamedArgs{
		"roleId": roleId,
	}

	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user role permissions: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.UserRolePermission])
}
