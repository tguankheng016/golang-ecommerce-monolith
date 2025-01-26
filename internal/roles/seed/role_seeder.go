package seed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/roles/constants"
	"github.com/tguankheng016/commerce-mono/internal/roles/models"
	"github.com/tguankheng016/commerce-mono/internal/roles/services"
)

type RoleSeeder struct {
	db *pgxpool.Pool
}

func NewRoleSeeder(db *pgxpool.Pool) RoleSeeder {
	return RoleSeeder{
		db: db,
	}
}

func (r RoleSeeder) SeedRoles(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
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

	roleManager := services.NewRoleManager(tx)

	roleFound, err := roleManager.GetRoleByName(ctx, constants.DefaultAdminRoleName)
	if err != nil {
		return err
	}

	if roleFound == nil {
		newAdminRole := &models.Role{
			Name:     constants.DefaultAdminRoleName,
			IsStatic: true,
		}

		if err := roleManager.CreateRole(ctx, newAdminRole); err != nil {
			return err
		}
	}

	roleFound, err = roleManager.GetRoleByName(ctx, constants.DefaultUserRoleName)

	if roleFound == nil {
		newUserRole := &models.Role{
			Name:      constants.DefaultUserRoleName,
			IsStatic:  true,
			IsDefault: true,
		}

		if err := roleManager.CreateRole(ctx, newUserRole); err != nil {
			return err
		}
	}

	return err
}
