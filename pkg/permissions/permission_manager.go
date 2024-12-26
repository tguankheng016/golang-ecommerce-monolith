package permissions

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"go.uber.org/zap"
)

type IPermissionManager interface {
	IsGranted(ctx context.Context, userId int64, permissionName string) (bool, error)
	GetGrantedPermissions(ctx context.Context, userId int64) (map[string]struct{}, error)
}

type IPermissionDbManager interface {
	GetGrantedPermissionsFromDb(ctx context.Context, userId int64) (map[string]struct{}, error)
}

type permissionManager struct {
	cacheManager *cache.Cache[string]
	dbManager    IPermissionDbManager
}

func NewPermissionManager(cacheManager *cache.Cache[string], dbManager IPermissionDbManager) IPermissionManager {
	return &permissionManager{
		cacheManager: cacheManager,
		dbManager:    dbManager,
	}
}

func ValidatePermissionName(grantedPermissions []string) error {
	allPermissions := GetAppPermissions().Items
	for _, permission := range grantedPermissions {
		if _, ok := allPermissions[permission]; !ok {
			return errors.New("invalid permission name")
		}
	}

	return nil
}

func (p *permissionManager) IsGranted(ctx context.Context, userId int64, permissionName string) (bool, error) {
	grantedPermissions, err := p.GetGrantedPermissions(ctx, userId)
	if err != nil {
		return false, err
	}

	_, ok := grantedPermissions[permissionName]
	return ok, nil
}

// GetGrantedPermissions gets all granted permissions for a given user id.
// It will first try to get the user roles and user permissions from cache.
// If the cache is not found, it will get the user roles and user permissions from database.
// For each user role, it will get the role permissions from cache.
// If the role permissions cache is not found, it will get the role permissions from database.
// It will return a map of permission names to empty struct.
// If there are any errors in the process, it will return an error.
func (p *permissionManager) GetGrantedPermissions(ctx context.Context, userId int64) (map[string]struct{}, error) {
	grantedPermissions := make(map[string]struct{})
	userProhibitedPermissions := make(map[string]struct{})

	// Check from user role cache first
	cachedUserRoles, err := p.cacheManager.Get(ctx, GenerateUserRoleCacheKey(userId))
	if err != nil {
		if !caching.CheckIsCacheValueNotFound(err) {
			logging.Logger.Error("Getting user roles caches err: ", zap.Error(err))
		}
		return p.dbManager.GetGrantedPermissionsFromDb(ctx, userId)
	}

	// Decode cache item from bytes
	var userRoleCacheItem UserRoleCacheItem
	if err := gob.NewDecoder(bytes.NewBuffer([]byte(cachedUserRoles))).Decode(&userRoleCacheItem); err != nil {
		return nil, err
	}

	// Get user permission caches
	cachedUserPermissions, err := p.cacheManager.Get(ctx, GenerateUserPermissionCacheKey(userId))
	if err != nil {
		if !caching.CheckIsCacheValueNotFound(err) {
			logging.Logger.Error("Getting user permissions caches err: ", zap.Error(err))
		}
		return p.dbManager.GetGrantedPermissionsFromDb(ctx, userId)
	}

	var userPermissionCacheItem UserPermissionCacheItem
	if err := gob.NewDecoder(bytes.NewBuffer([]byte(cachedUserPermissions))).Decode(&userPermissionCacheItem); err != nil {
		return nil, err
	}

	if len(userPermissionCacheItem.Permissions) > 0 {
		grantedPermissions = userPermissionCacheItem.Permissions
	}

	if len(userPermissionCacheItem.ProhibitedPermissions) > 0 {
		userProhibitedPermissions = userPermissionCacheItem.ProhibitedPermissions
	}

	for _, r := range userRoleCacheItem.RoleIds {
		cachedRolePermissions, err := p.cacheManager.Get(ctx, GenerateRolePermissionCacheKey(r))
		if err != nil {
			if !caching.CheckIsCacheValueNotFound(err) {
				logging.Logger.Error("Getting role permissions caches err: ", zap.Error(err))
			}
			return p.dbManager.GetGrantedPermissionsFromDb(ctx, userId)
		}

		var rolePermissionCacheItem RolePermissionCacheItem
		if err := gob.NewDecoder(bytes.NewBuffer([]byte(cachedRolePermissions))).Decode(&rolePermissionCacheItem); err != nil {
			return nil, err
		}

		// add those unique and non prohibited (user level) role permission
		for p := range rolePermissionCacheItem.Permissions {
			if _, ok := userProhibitedPermissions[p]; !ok {
				if _, ok := grantedPermissions[p]; !ok {
					grantedPermissions[p] = struct{}{}
				}
			}
		}
	}
	return grantedPermissions, nil
}
