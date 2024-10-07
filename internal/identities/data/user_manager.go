package data

import (
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/security"
	"gorm.io/gorm"
)

type IUserManager interface {
	CreateUser(user *models.User, password string) error
	UpdateUser(user *models.User) error
	//DeleteUser(user *models.User) error
}

type userManager struct {
	db     *gorm.DB
	logger logger.ILogger
}

func NewUserManager(db *gorm.DB, logger logger.ILogger) IUserManager {
	return &userManager{
		db:     db,
		logger: logger,
	}
}

func (u *userManager) CreateUser(user *models.User, password string) error {
	if err := u.ValidateUser(user); err != nil {
		return err
	}

	if password != "" {
		hashPassword, err := security.HashPassword(password)
		if err != nil {
			return err
		}
		user.Password = hashPassword
	}

	user.NormalizedUserName = strings.ToUpper(user.UserName)
	user.NormalizedEmail = strings.ToUpper(user.Email)

	securityStamp, err := uuid.NewV4()
	if err != nil {
		return err
	}

	user.SecurityStamp = securityStamp

	if err := u.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userManager) UpdateUser(user *models.User) error {
	return nil
}

func (u *userManager) ValidateUser(user *models.User) error {
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
