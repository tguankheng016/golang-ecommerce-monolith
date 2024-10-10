package data

import (
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/core/helpers"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/security"
	"gorm.io/gorm"
)

type IUserManager interface {
	CreateUser(user *models.User, password string) error
	UpdateUser(user *models.User, password string) error
	GetUserRoles(user *models.User) ([]int64, error)
	UpdateUserRoles(user *models.User, roles []int64) error
	AddToRoles(user *models.User, roles []int64) error
	RemoveToRoles(user *models.User, roles []int64) error
}

type userManager struct {
	db *gorm.DB
}

func NewUserManager(db *gorm.DB) IUserManager {
	return &userManager{
		db: db,
	}
}

func (u *userManager) CreateUser(user *models.User, password string) error {
	if err := u.validateUser(user); err != nil {
		return err
	}

	if password == "" {
		return errors.New("Password is required")
	}

	u.hashUserPassword(user, password)

	user.NormalizedUserName = strings.ToUpper(user.UserName)
	user.NormalizedEmail = strings.ToUpper(user.Email)

	securityStamp, err := uuid.NewV4()
	if err != nil {
		return err
	}

	user.SecurityStamp = securityStamp

	if err := u.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userManager) UpdateUser(user *models.User, password string) error {
	if err := u.validateUser(user); err != nil {
		return err
	}

	if password != "" {
		u.hashUserPassword(user, password)
	}

	user.NormalizedUserName = strings.ToUpper(user.UserName)
	user.NormalizedEmail = strings.ToUpper(user.Email)

	if err := u.db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userManager) GetUserRoles(user *models.User) ([]int64, error) {
	userRoleIds := make([]int64, 0)
	if err := u.db.Model(&models.User{}).
		Where("id = ?", user.Id).
		Select("user_roles.role_id").
		Joins("join user_roles on user_roles.user_id = users.id").
		Scan(&userRoleIds).Error; err != nil {

		return nil, err
	}

	return userRoleIds, nil
}

func (u *userManager) UpdateUserRoles(user *models.User, roles []int64) error {
	userRoleIds, err := u.GetUserRoles(user)
	if err != nil {
		return err
	}

	var roleIdsToAdd []int64
	for _, roleId := range roles {
		if !helpers.SliceContains(userRoleIds, roleId) {
			roleIdsToAdd = append(roleIdsToAdd, roleId)
		}
	}

	var roleIdsToRemove []int64
	for _, userRoleId := range userRoleIds {
		if !helpers.SliceContains(roles, userRoleId) {
			roleIdsToRemove = append(roleIdsToRemove, userRoleId)
		}
	}

	if len(roleIdsToAdd) > 0 {
		if err := u.AddToRoles(user, roleIdsToAdd); err != nil {
			return err
		}
	}

	if len(roleIdsToRemove) > 0 {
		if err := u.RemoveToRoles(user, roleIdsToRemove); err != nil {
			return err
		}
	}

	return nil
}

func (u *userManager) AddToRoles(user *models.User, roles []int64) error {
	for _, roleId := range roles {
		var role models.Role
		if err := u.db.First(&role, roleId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return err
		}

		var count int64
		if err := u.db.Model(&models.User{}).
			Where("id = ?", user.Id).
			Select("user_roles.role_id").
			Joins("join user_roles on user_roles.user_id = users.id").
			Count(&count).Error; err != nil {

			return err
		}

		if count == 0 {
			if err := u.db.Model(user).Association("Roles").Append(&role); err != nil {
				return errors.Wrap(err, "error in the assigning admin role")
			}
		}
	}
	return nil
}

func (u *userManager) RemoveToRoles(user *models.User, roles []int64) error {
	for _, roleId := range roles {
		var role models.Role
		if err := u.db.First(&role, roleId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return err
		}

		if err := u.db.Model(user).Association("Roles").Delete(&models.Role{Id: roleId}); err != nil {
			return err
		}
	}

	return nil
}

func (u *userManager) validateUser(user *models.User) error {
	var count int64
	if err := u.db.Model(&models.User{}).Where("id <> ? AND normalized_user_name = ?", user.Id, strings.ToUpper(user.UserName)).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Username has been taken")
	}
	if err := u.db.Model(&models.User{}).Where("id <> ? AND normalized_email = ?", user.Id, strings.ToUpper(user.Email)).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Email has been taken")
	}

	return nil
}

func (u *userManager) hashUserPassword(user *models.User, password string) error {
	if len(password) < 6 {
		return errors.New("Password must be at least 8 characters long")
	}

	hashPassword, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashPassword

	return nil
}
