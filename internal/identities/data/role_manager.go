package data

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"gorm.io/gorm"
)

type IRoleManager interface {
	CreateRole(role *models.Role) error
	UpdateRole(role *models.Role) error
}

type roleManager struct {
	db *gorm.DB
}

func NewRoleManager(db *gorm.DB) IRoleManager {
	return &roleManager{
		db: db,
	}
}

func (r *roleManager) CreateRole(role *models.Role) error {
	if err := r.validateRole(role); err != nil {
		return err
	}

	if err := r.db.Create(role).Error; err != nil {
		return err
	}

	return nil
}

func (r *roleManager) UpdateRole(role *models.Role) error {
	if err := r.validateRole(role); err != nil {
		return err
	}

	if err := r.db.Save(role).Error; err != nil {
		return err
	}

	return nil
}

func (r *roleManager) validateRole(role *models.Role) error {
	var count int64
	if err := r.db.Model(&models.Role{}).Where("id <> ? AND UPPER(name) = ?", role.Id, strings.ToUpper(role.Name)).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Email has been taken")
	}

	return nil
}
