package permissions

import (
	"bytes"
	"context"
	"strings"
	"time"

	"encoding/gob"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	identityConsts "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"gorm.io/gorm"
)

type IPermissionManager interface {
	IsGranted(ctx context.Context, userId int64, permissionName string) (bool, error)
	SetUserPermissions(ctx context.Context, userId int64) (map[string]struct{}, error)
	SetRolePermissions(ctx context.Context, roleId int64) (map[string]struct{}, error)
	RemoveUserRoleCaches(ctx context.Context, userId int64)
	RemoveUserPermissionCaches(ctx context.Context, userId int64)
	RemoveRolePermissionCaches(cctx context.Context, roleId int64)
}

type permissionManager struct {
	db     *gorm.DB
	client *redis.Client
	logger logger.ILogger
}

const (
	DefaultCacheExpiration = 1 * time.Hour
)

func NewPermissionManager(db *gorm.DB, client *redis.Client, logger logger.ILogger) IPermissionManager {
	return &permissionManager{
		db:     db,
		client: client,
		logger: logger,
	}
}

func (p *permissionManager) IsGranted(ctx context.Context, userId int64, permissionName string) (bool, error) {
	// Check from cache first
	cachedUserRoles, err := p.client.Get(ctx, GenerateUserRoleCacheKey(userId)).Result()
	if err != nil {
		if err != redis.Nil {
			p.logger.Error(err)
		}
		return p.isGrantedFromDb(ctx, userId, permissionName)
	}

	var userRoleCacheItem UserRoleCacheItem
	if err := gob.NewDecoder(bytes.NewBuffer([]byte(cachedUserRoles))).Decode(&userRoleCacheItem); err != nil {
		return false, err
	}

	cachedUserPermissions, err := p.client.Get(ctx, GenerateUserPermissionCacheKey(userId)).Result()
	if err != nil && err != redis.Nil {
		p.logger.Error(err)
	}

	var userPermissionCacheItem UserPermissionCacheItem
	if err := gob.NewDecoder(bytes.NewBuffer([]byte(cachedUserPermissions))).Decode(&userPermissionCacheItem); err != nil {
		if err != redis.Nil {
			p.logger.Error(err)
		}
	} else {
		if _, ok := userPermissionCacheItem.Permissions[permissionName]; ok {
			return true, nil
		}

		if _, ok := userPermissionCacheItem.ProhibitedPermissions[permissionName]; ok {
			return false, nil
		}
	}

	for _, r := range userRoleCacheItem.RoleIds {
		cachedRolePermissions, err := p.client.Get(ctx, GenerateRolePermissionCacheKey(r)).Result()
		if err != nil {
			if err != redis.Nil {
				p.logger.Error(err)
			}
			continue
		}

		var rolePermissionCacheItem RolePermissionCacheItem
		if err := gob.NewDecoder(bytes.NewBuffer([]byte(cachedRolePermissions))).Decode(&rolePermissionCacheItem); err != nil {
			p.logger.Error(err)
			continue
		}

		if _, ok := rolePermissionCacheItem.Permissions[permissionName]; ok {
			return true, nil
		}
	}

	return false, nil
}

func (p *permissionManager) isGrantedFromDb(ctx context.Context, userId int64, permissionName string) (bool, error) {
	grantedPermissions, err := p.SetUserPermissions(ctx, userId)
	if err != nil {
		return false, err
	}

	_, ok := grantedPermissions[permissionName]

	return ok, nil
}

func (p *permissionManager) SetUserPermissions(ctx context.Context, userId int64) (map[string]struct{}, error) {
	var user models.User
	if err := p.db.Model(&models.User{}).Preload("Roles").First(&user, userId).Error; err != nil {
		return nil, err
	}

	// Cache User Role Ids
	roleIds := make([]int64, 0)
	for _, role := range user.Roles {
		roleIds = append(roleIds, role.Id)
	}

	userRoleCacheItem, err := marshalCacheItem(NewUserRoleCacheItem(userId, roleIds))
	if err != nil {
		return nil, err
	}
	if err := p.client.Set(ctx, GenerateUserRoleCacheKey(userId), userRoleCacheItem, DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		p.logger.Error(err)
	}

	var userPermissions []models.UserRolePermission

	if err := p.db.Model(&models.UserRolePermission{}).Where("user_id = ?", userId).Find(&userPermissions).Error; err != nil {
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

	userPermissionCacheItem, err := marshalCacheItem(NewUserPermissionCacheItem(userId, grantedUserPermissions, prohibitedUserPermissions))
	if err != nil {
		return nil, err
	}
	if err := p.client.Set(ctx, GenerateUserPermissionCacheKey(userId), userPermissionCacheItem, DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		p.logger.Error(err)
	}

	for _, role := range user.Roles {
		rolePermissions, err := p.SetRolePermissions(ctx, role.Id)
		if err != nil {
			return nil, err
		}

		for permission, _ := range rolePermissions {
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

func (p *permissionManager) SetRolePermissions(ctx context.Context, roleId int64) (map[string]struct{}, error) {
	var role models.Role
	if err := p.db.Model(&models.Role{}).First(&role, roleId).Error; err != nil {
		return nil, err
	}

	isAdmin := strings.EqualFold(role.Name, identityConsts.DefaultAdminRoleName)

	var rolePermissions []models.UserRolePermission

	if err := p.db.Model(&models.UserRolePermission{}).Where("role_id = ?", roleId).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	if isAdmin {
		allPermissions := GetAppPermissions().Items

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

		rolePermissionCacheItem, err := marshalCacheItem(NewRolePermissionCacheItem(roleId, grantedPermissions))
		if err != nil {
			return nil, err
		}
		if err := p.client.Set(ctx, GenerateRolePermissionCacheKey(roleId), rolePermissionCacheItem, DefaultCacheExpiration).Err(); err != nil {
			// Dont return just log
			p.logger.Error(err)
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

func (p *permissionManager) RemoveUserRoleCaches(ctx context.Context, userId int64) {
	if err := p.client.Del(ctx, GenerateUserRoleCacheKey(userId)).Err(); err != nil {
		p.logger.Error(err)
	}
}

func (p *permissionManager) RemoveUserPermissionCaches(ctx context.Context, userId int64) {
	if err := p.client.Del(ctx, GenerateUserPermissionCacheKey(userId)).Err(); err != nil {
		p.logger.Error(err)
	}
}

func (p *permissionManager) RemoveRolePermissionCaches(cctx context.Context, roleId int64) {
	if err := p.client.Del(cctx, GenerateRolePermissionCacheKey(roleId)).Err(); err != nil {
		p.logger.Error(err)
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

func marshalCacheItem(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
