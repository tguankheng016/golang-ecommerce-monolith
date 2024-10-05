package data

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/security"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

func DataSeeder(gorm *gorm.DB) error {
	if err := seedRole(gorm); err != nil {
		return err
	}

	if err := seedUser(gorm); err != nil {
		return err
	}

	return nil
}

func seedRole(gorm *gorm.DB) error {
	if (gorm.Find(&models.Role{}).RowsAffected <= 0) {
		adminRole := &models.Role{
			Name:      constants.DefaultAdminRoleName,
			CreatedAt: time.Now(),
		}

		if err := gorm.Create(adminRole).Error; err != nil {
			return errors.Wrap(err, "error in the inserting role into the database.")
		}
	}

	return nil
}

func seedUser(gorm *gorm.DB) error {
	if (gorm.Find(&models.User{}).RowsAffected <= 0) {
		pass, err := security.HashPassword("123qwe")
		if err != nil {
			return err
		}

		securityStamp, err := uuid.NewV4()
		if err != nil {
			return err
		}

		adminUser := &models.User{
			FirstName:     "admin",
			LastName:      "Tan",
			UserName:      constants.DefaultAdminUsername,
			Email:         "admin@testgk.com",
			SecurityStamp: securityStamp,
			Password:      pass,
			CreatedAt:     time.Now(),
		}

		if err := gorm.Create(adminUser).Error; err != nil {
			return errors.Wrap(err, "error in the inserting admin user into the database.")
		}

		var adminRole models.Role

		if err := gorm.Where("name = ?", constants.DefaultAdminRoleName).First(&adminRole).Error; err != nil {
			return errors.Wrap(err, "error in the selecting default admin role")
		}

		if err := gorm.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
			return errors.Wrap(err, "error in the assigning admin role")
		}

		normalUser := &models.User{
			FirstName:     "User",
			LastName:      "Tan",
			UserName:      "gkuser123",
			Email:         "user@testgk.com",
			SecurityStamp: securityStamp,
			Password:      pass,
			CreatedAt:     time.Now(),
		}

		if err := gorm.Create(normalUser).Error; err != nil {
			return errors.Wrap(err, "error in the inserting normal user into the database.")
		}
	}

	return nil
}
