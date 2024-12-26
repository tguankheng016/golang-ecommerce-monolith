package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tguankheng016/commerce-mono/internal/roles/models"
	httpServer "github.com/tguankheng016/commerce-mono/pkg/http"
)

func (r roleManager) CreateRole(ctx context.Context, role *models.Role) error {
	if err := r.validateRoleName(ctx, role); err != nil {
		return err
	}

	query := `
		INSERT INTO roles (name, normalized_name, is_default, is_static, created_by, is_deleted) 
		VALUES (@name, @normalized_name, @is_default, @is_static, @created_by, false)
		RETURNING id;
	`

	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		role.CreatedBy.Int64 = currentUserId
		role.CreatedBy.Valid = true
	}

	args := pgx.NamedArgs{
		"name":            role.Name,
		"normalized_name": strings.ToUpper(role.Name),
		"is_default":      role.IsDefault,
		"is_static":       role.IsStatic,
		"created_by":      role.CreatedBy,
	}

	// Variable to store the returned ID
	var insertedID int64

	// Execute the insert query and retrieve the inserted ID
	if err := r.db.QueryRow(context.Background(), query, args).Scan(&insertedID); err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	role.Id = insertedID

	return nil
}

func (r roleManager) CreateRolePermission(ctx context.Context, roleId int64, permission string, isGranted bool) error {
	query := `
		INSERT INTO user_role_permissions (
			name,
			role_id, 
			is_granted
		) 
		VALUES (
			@name,
			@role_id, 
			@is_granted
		)
	`

	args := pgx.NamedArgs{
		"name":       permission,
		"role_id":    roleId,
		"is_granted": isGranted,
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to insert role permission: %w", err)
	}

	return nil
}

func (r roleManager) UpdateRole(ctx context.Context, role *models.Role) error {
	if err := r.validateRoleName(ctx, role); err != nil {
		return err
	}

	query := `
		UPDATE roles
		SET 
			name = @name, 
			normalized_name = @normalized_name, 
			is_default = @is_default, 
			updated_at = @updated_at,
			updated_by = @updated_by
		WHERE 
			id = @id
	`

	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		role.UpdatedBy.Int64 = currentUserId
		role.UpdatedBy.Valid = true
	}

	args := pgx.NamedArgs{
		"name":            role.Name,
		"normalized_name": strings.ToUpper(role.Name),
		"is_default":      role.IsDefault,
		"id":              role.Id,
		"updated_at":      time.Now(),
		"updated_by":      role.UpdatedBy,
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to update role: %w", err)
	}

	return nil
}

func (r roleManager) DeleteRole(ctx context.Context, roleId int64) error {
	query := "DELETE FROM user_role_permissions WHERE role_id = $1"
	if _, err := r.db.Exec(ctx, query, roleId); err != nil {
		return fmt.Errorf("unable to delete role permission: %w", err)
	}

	query = `
		UPDATE roles
		SET 
			is_deleted = true,
			deleted_at = @deleted_at,
			deleted_by = @deleted_by
		WHERE 
			id = @id
	`

	deletedUserId := &sql.NullInt64{}

	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		deletedUserId.Int64 = currentUserId
		deletedUserId.Valid = true
	}

	args := pgx.NamedArgs{
		"id":         roleId,
		"deleted_at": time.Now(),
		"deleted_by": deletedUserId,
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to delete role: %w", err)
	}

	return nil
}

func (r roleManager) DeleteRolePermission(ctx context.Context, roleId int64, permission string) error {
	query := "DELETE FROM user_role_permissions WHERE role_id = $1 AND name = $2"

	if _, err := r.db.Exec(ctx, query, roleId, permission); err != nil {
		return fmt.Errorf("unable to delete role permission: %w", err)
	}

	return nil
}

func (r roleManager) validateRoleName(ctx context.Context, role *models.Role) error {
	query := "select count(*) from roles where normalized_name = @roleName and is_deleted = false and id != @id"

	args := pgx.NamedArgs{
		"roleName": strings.ToUpper(role.Name),
		"id":       role.Id,
	}

	var count int

	if err := r.db.QueryRow(ctx, query, args).Scan(&count); err != nil {
		return fmt.Errorf("unable to count role: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("role name %s is taken", role.Name)
	}

	return nil
}
