package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tguankheng016/commerce-mono/internal/users/models"
	httpServer "github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/security"
)

func (u userManager) CreateUser(ctx context.Context, user *models.User, password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if err := u.validateUserName(ctx, user); err != nil {
		return err
	}

	if err := u.validateUserEmail(ctx, user); err != nil {
		return err
	}

	if err := hashUserPassword(user, password); err != nil {
		return err
	}

	securityStamp, err := uuid.NewV6()
	if err != nil {
		return err
	}

	user.SecurityStamp = securityStamp

	query := `
		INSERT INTO users (
			first_name, 
			last_name, 
			user_name, 
			normalized_user_name, 
			email, 
			normalized_email, 
			password_hash,
			security_stamp,
			created_by,
			is_deleted
		) 
		VALUES (
			@first_name, 
			@last_name, 
			@user_name, 
			@normalized_user_name, 
			@email, 
			@normalized_email, 
			@password_hash,
			@security_stamp,
			@created_by,
			false
		)
		RETURNING id;
	`
	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		user.CreatedBy.Int64 = currentUserId
		user.CreatedBy.Valid = true
	}

	args := pgx.NamedArgs{
		"first_name":           user.FirstName,
		"last_name":            user.LastName,
		"user_name":            user.UserName,
		"normalized_user_name": strings.ToUpper(user.UserName),
		"email":                user.Email,
		"normalized_email":     strings.ToUpper(user.Email),
		"password_hash":        user.PasswordHash,
		"security_stamp":       user.SecurityStamp,
		"created_by":           user.CreatedBy,
	}

	// Variable to store the returned ID
	var insertedID int64

	// Execute the insert query and retrieve the inserted ID
	if err = u.db.QueryRow(ctx, query, args).Scan(&insertedID); err != nil {
		return fmt.Errorf("unable to insert user: %w", err)
	}

	user.Id = insertedID

	return nil
}

func (u userManager) CreateUserRole(ctx context.Context, userId int64, roleId int64) error {
	query := `
		INSERT INTO user_roles (
			user_id, 
			role_id
		) 
		VALUES (
			@user_id, 
			@role_id
		)
	`

	args := pgx.NamedArgs{
		"user_id": userId,
		"role_id": roleId,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to insert user role: %w", err)
	}

	return nil
}

func (u userManager) CreateUserPermission(ctx context.Context, userId int64, permission string, isGranted bool) error {
	query := `
		INSERT INTO user_role_permissions (
			name,
			user_id, 
			is_granted
		) 
		VALUES (
			@name,
			@user_id, 
			@is_granted
		)
	`

	args := pgx.NamedArgs{
		"name":       permission,
		"user_id":    userId,
		"is_granted": isGranted,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to insert user permission: %w", err)
	}

	return nil
}

func (u userManager) UpdateUser(ctx context.Context, user *models.User, password string) error {
	if err := u.validateUserName(ctx, user); err != nil {
		return err
	}

	if err := u.validateUserEmail(ctx, user); err != nil {
		return err
	}

	if password != "" {
		if err := hashUserPassword(user, password); err != nil {
			return err
		}
	}

	query := `
		UPDATE users
		SET 
			first_name = @first_name, 
			last_name = @last_name, 
			user_name = @user_name, 
			normalized_user_name = @normalized_user_name, 
			email = @email, 
			normalized_email = @normalized_email, 
			password_hash = @password_hash,
			security_stamp = @security_stamp,
			updated_at = @updated_at,
			updated_by = @updated_by
		WHERE 
			id = @id
	`
	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		user.UpdatedBy.Int64 = currentUserId
		user.UpdatedBy.Valid = true
	}

	args := pgx.NamedArgs{
		"first_name":           user.FirstName,
		"last_name":            user.LastName,
		"user_name":            user.UserName,
		"normalized_user_name": strings.ToUpper(user.UserName),
		"email":                user.Email,
		"normalized_email":     strings.ToUpper(user.Email),
		"password_hash":        user.PasswordHash,
		"security_stamp":       user.SecurityStamp,
		"id":                   user.Id,
		"updated_at":           time.Now(),
		"updated_by":           user.UpdatedBy,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to update user: %w", err)
	}

	return nil
}

func (u userManager) UpdateUserRoles(ctx context.Context, user *models.User, roles []int64) (bool, error) {
	userRoleIds, err := u.GetUserRoleIds(ctx, user.Id)
	if err != nil {
		return false, err
	}

	var roleIdsToAdd []int64
	for _, roleId := range roles {
		if !slices.Contains(userRoleIds, roleId) && roleId != 0 {
			roleIdsToAdd = append(roleIdsToAdd, roleId)
		}
	}

	var roleIdsToRemove []int64
	for _, userRoleId := range userRoleIds {
		if !slices.Contains(roles, userRoleId) {
			roleIdsToRemove = append(roleIdsToRemove, userRoleId)
		}
	}

	if len(roleIdsToAdd) > 0 {
		for _, roleId := range roleIdsToAdd {
			if err := u.CreateUserRole(ctx, user.Id, roleId); err != nil {
				return false, err
			}
		}
	}

	if len(roleIdsToRemove) > 0 {
		query := "DELETE FROM user_roles WHERE user_id = $1 and role_id = ANY($2)"
		if _, err := u.db.Exec(ctx, query, user.Id, roleIdsToRemove); err != nil {
			return false, fmt.Errorf("unable to delete user role: %w", err)
		}
	}

	return len(roleIdsToRemove) > 0 || len(roleIdsToAdd) > 0, nil
}

func (u userManager) DeleteUser(ctx context.Context, userId int64) error {
	query := "DELETE FROM user_roles WHERE user_id = $1"
	if _, err := u.db.Exec(ctx, query, userId); err != nil {
		return fmt.Errorf("unable to delete user roles: %w", err)
	}

	query = "DELETE FROM user_role_permissions WHERE user_id = $1"
	if _, err := u.db.Exec(ctx, query, userId); err != nil {
		return fmt.Errorf("unable to delete user permission: %w", err)
	}

	query = `
		UPDATE users
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
		"id":         userId,
		"deleted_at": time.Now(),
		"deleted_by": deletedUserId,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}

	return nil
}

func (u userManager) DeleteUserRole(ctx context.Context, userId int64, roleId int64) error {
	query := "DELETE FROM user_roles WHERE user_id = $1 and role_id = $2"
	if _, err := u.db.Exec(ctx, query, userId, roleId); err != nil {
		return fmt.Errorf("unable to delete user role: %w", err)
	}

	return nil
}

func (u userManager) DeleteUserPermission(ctx context.Context, userId int64, permission string) error {
	query := "DELETE FROM user_role_permissions WHERE user_id = $1 and name = $2"
	if _, err := u.db.Exec(ctx, query, userId, permission); err != nil {
		return fmt.Errorf("unable to delete user permission: %w", err)
	}

	return nil
}

func (u userManager) DeleteUserPermissions(ctx context.Context, userId int64) error {
	query := "DELETE FROM user_role_permissions WHERE user_id = $1"
	if _, err := u.db.Exec(ctx, query, userId); err != nil {
		return fmt.Errorf("unable to delete user permissions: %w", err)
	}

	return nil
}

func (u userManager) validateUserName(ctx context.Context, user *models.User) error {
	query := "select count(*) from users where normalized_user_name = @username and is_deleted = false and id != @id"

	args := pgx.NamedArgs{
		"username": strings.ToUpper(user.UserName),
		"id":       user.Id,
	}

	var count int

	if err := u.db.QueryRow(ctx, query, args).Scan(&count); err != nil {
		return fmt.Errorf("unable to count user: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("username %s is taken", user.UserName)
	}

	return nil
}

func (u userManager) validateUserEmail(ctx context.Context, user *models.User) error {
	query := "select count(*) from users where normalized_email = @email and is_deleted = false and id != @id"

	args := pgx.NamedArgs{
		"email": strings.ToUpper(user.Email),
		"id":    user.Id,
	}

	var count int

	if err := u.db.QueryRow(ctx, query, args).Scan(&count); err != nil {
		return fmt.Errorf("unable to count user: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("email %s is taken", user.Email)
	}

	return nil
}

func hashUserPassword(user *models.User, password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	hashPassword, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = hashPassword

	return nil
}
