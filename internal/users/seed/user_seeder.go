package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	roleConsts "github.com/tguankheng016/commerce-mono/internal/roles/constants"
	roleService "github.com/tguankheng016/commerce-mono/internal/roles/services"
	"github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/internal/users/models"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
)

type UserSeeder struct {
	db *pgxpool.Pool
}

func NewUserSeeder(db *pgxpool.Pool) UserSeeder {
	return UserSeeder{
		db: db,
	}
}

func (u UserSeeder) SeedUsers(ctx context.Context) error {
	tx, err := u.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			// Rollback the transaction in case of error
			tx.Rollback(ctx)
		} else {
			// Commit the transaction if no error occurs
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("unable to commit transaction: %w", err)
			}
		}
	}()

	userManager := userService.NewUserManager(tx)
	roleManager := roleService.NewRoleManager(tx)

	adminUserFound, err := userManager.GetUserByUserName(ctx, constants.DefaultAdminUserName)
	if err != nil {
		return err
	}

	if adminUserFound == nil {
		adminRole, err := roleManager.GetRoleByName(ctx, roleConsts.DefaultAdminRoleName)
		if err != nil {
			return err
		}
		if adminRole == nil {
			return errors.New("admin role not found")
		}

		newAdminUser := &models.User{
			UserName:  constants.DefaultAdminUserName,
			FirstName: constants.DefaultAdminUserName,
			LastName:  "Tan",
			Email:     "admin@testgk.com",
		}

		if err := userManager.CreateUser(ctx, newAdminUser, "123qwe"); err != nil {
			return err
		}

		if err := userManager.CreateUserRole(ctx, newAdminUser.Id, adminRole.Id); err != nil {
			return err
		}

		normalUserFound, err := userManager.GetUserByUserName(ctx, constants.DefaultUserUserName)
		if err != nil {
			return err
		}

		if normalUserFound == nil {
			userRole, err := roleManager.GetRoleByName(ctx, roleConsts.DefaultUserRoleName)
			if err != nil {
				return err
			}
			if userRole == nil {
				return errors.New("user role not found")
			}

			newNormalUser := &models.User{
				UserName:  constants.DefaultUserUserName,
				FirstName: constants.DefaultUserUserName,
				LastName:  "Tan",
				Email:     "gktan@testgk.com",
			}

			if err := userManager.CreateUser(ctx, newNormalUser, "123qwe"); err != nil {
				return err
			}

			if err := userManager.CreateUserRole(ctx, newNormalUser.Id, userRole.Id); err != nil {
				return err
			}
		}
	}

	return err
}
