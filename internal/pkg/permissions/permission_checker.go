package permissions

import (
	"bytes"
	"context"
	"strings"
	"time"

	"encoding/gob"

	"github.com/redis/go-redis/v9"
	identityConsts "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"gorm.io/gorm"
)

type IPermissionChecker interface {
	IsGranted(userId int64, permissionName string) (bool, error)
}

type permissionChecker struct {
	db     *gorm.DB
	client *redis.Client
	logger logger.ILogger
}

const (
	DefaultCacheExpiration = 1 * time.Hour
)

func NewPermissionChecker(db *gorm.DB, client *redis.Client, logger logger.ILogger) IPermissionChecker {
	return &permissionChecker{
		db:     db,
		client: client,
		logger: logger,
	}
}

func (p *permissionChecker) IsGranted(userId int64, permissionName string) (bool, error) {
	// Check from cache first
	cachedUserRoles, err := p.client.Get(context.Background(), GenerateUserRoleCacheKey(userId)).Result()
	if err != nil {
		if err != redis.Nil {
			p.logger.Error(err)
		}
		return p.isGrantedFromDb(userId, permissionName)
	}

	var userRoleCacheItem UserRoleCacheItem
	if err := gob.NewDecoder(bytes.NewBuffer([]byte(cachedUserRoles))).Decode(&userRoleCacheItem); err != nil {
		return false, err
	}

	cachedUserPermissions, err := p.client.Get(context.Background(), GenerateUserPermissionCacheKey(userId)).Result()
	if err != nil && err != redis.Nil {
		p.logger.Error(err)
	}

	var userPermissionCacheItem UserPermissionCacheItem
	if err := gob.NewDecoder(bytes.NewBuffer([]byte(cachedUserPermissions))).Decode(&userPermissionCacheItem); err != nil {
		p.logger.Error(err)
	} else {
		if _, ok := userPermissionCacheItem.Permissions[permissionName]; ok {
			return true, nil
		}

		if _, ok := userPermissionCacheItem.ProhibitedPermissions[permissionName]; ok {
			return false, nil
		}
	}

	for _, r := range userRoleCacheItem.RoleIds {
		cachedRolePermissions, err := p.client.Get(context.Background(), GenerateRolePermissionCacheKey(r)).Result()
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

func (p *permissionChecker) isGrantedFromDb(userId int64, permissionName string) (bool, error) {
	grantedPermissions, err := p.setUserPermissions(userId)
	if err != nil {
		return false, err
	}

	_, ok := grantedPermissions[permissionName]

	return ok, nil
}

func (p *permissionChecker) setUserPermissions(userId int64) (map[string]string, error) {
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
	if err := p.client.Set(context.Background(), GenerateUserRoleCacheKey(userId), userRoleCacheItem, DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		p.logger.Error(err)
	}

	var userPermissions []models.UserRolePermission

	if err := p.db.Model(&models.UserRolePermission{}).Where("user_id = ?", userId).Find(&userPermissions).Error; err != nil {
		return nil, err
	}

	grantedPermissions := make(map[string]string)
	grantedUserPermissions := make(map[string]string)
	prohibitedUserPermissions := make(map[string]string)
	for _, permission := range userPermissions {
		if permission.IsGranted {
			grantedUserPermissions[permission.Name] = permission.Name
			grantedPermissions[permission.Name] = permission.Name
		} else {
			prohibitedUserPermissions[permission.Name] = permission.Name
		}
	}

	userPermissionCacheItem, err := marshalCacheItem(NewUserPermissionCacheItem(userId, grantedUserPermissions, prohibitedUserPermissions))
	if err != nil {
		return nil, err
	}
	if err := p.client.Set(context.Background(), GenerateUserPermissionCacheKey(userId), userPermissionCacheItem, DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		p.logger.Error(err)
	}

	for _, role := range user.Roles {
		rolePermissions, err := p.setRolePermissions(role.Id)
		if err != nil {
			return nil, err
		}

		for _, permission := range rolePermissions {
			if _, ok := prohibitedUserPermissions[permission]; !ok {
				if _, ok := grantedPermissions[permission]; !ok {
					// key does not exists
					grantedPermissions[permission] = permission
				}
			}
		}
	}

	return grantedPermissions, nil
}

func (p *permissionChecker) setRolePermissions(roleId int64) (map[string]string, error) {
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
		prohibitedPermissions := make(map[string]string)
		for _, permission := range rolePermissions {
			if !permission.IsGranted {
				prohibitedPermissions[permission.Name] = permission.Name
			}
		}

		// Excluded prohibited permissions
		grantedPermissions := make(map[string]string)
		for key := range allPermissions {
			if _, ok := prohibitedPermissions[key]; !ok {
				grantedPermissions[key] = key
			}
		}

		rolePermissionCacheItem, err := marshalCacheItem(NewRolePermissionCacheItem(roleId, grantedPermissions))
		if err != nil {
			return nil, err
		}
		if err := p.client.Set(context.Background(), GenerateRolePermissionCacheKey(roleId), rolePermissionCacheItem, DefaultCacheExpiration).Err(); err != nil {
			// Dont return just log
			p.logger.Error(err)
		}

		return grantedPermissions, nil
	} else {
		grantedPermissions := make(map[string]string)
		for _, permission := range rolePermissions {
			if permission.IsGranted {
				grantedPermissions[permission.Name] = permission.Name
			}
		}

		return grantedPermissions, nil
	}
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
