package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/tguankheng016/commerce-mono/internal/roles/models"
	userModel "github.com/tguankheng016/commerce-mono/internal/users/models"
	"github.com/tguankheng016/commerce-mono/pkg/core/pagination"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

type IRoleManager interface {
	GetRoles(ctx context.Context, pageRequest *pagination.PageRequest) ([]models.Role, int, error)
	GetRolesCount(ctx context.Context) (int, error)
	GetRoleById(ctx context.Context, roleId int64) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	GetRolePermission(ctx context.Context, roleId int64, permission string) (*userModel.UserRolePermission, error)
	GetRolePermissions(ctx context.Context, roleId int64, isGranted *bool) ([]userModel.UserRolePermission, error)

	CreateRole(ctx context.Context, role *models.Role) error
	CreateRolePermission(ctx context.Context, roleId int64, permission string, isGranted bool) error

	UpdateRole(ctx context.Context, role *models.Role) error

	DeleteRole(ctx context.Context, roleId int64) error
	DeleteRolePermission(ctx context.Context, roleId int64, permission string) error
}

type roleManager struct {
	db postgres.IPgxDbConn
}

func NewRoleManager(db postgres.IPgxDbConn) IRoleManager {
	return roleManager{
		db: db,
	}
}

func (r roleManager) GetRoles(ctx context.Context, pageRequest *pagination.PageRequest) ([]models.Role, int, error) {
	query := `SELECT %s FROM roles WHERE is_deleted = false %s %s %s`
	whereExpr := ""
	sortExpr := ""
	paginateExpr := ""
	count := 0

	args := pgx.NamedArgs{}

	if pageRequest != nil {
		if pageRequest.Filters != "" {
			whereExpr = `
				AND (
					normalized_name like @filters
				)
			`

			args["filters"] = fmt.Sprintf("%%%s%%", strings.ToUpper(pageRequest.Filters))
		}

		if pageRequest.Sorting != "" {
			sortingFields := []string{"name"}
			if err := pageRequest.SanitizeSorting(sortingFields...); err != nil {
				return nil, 0, err
			}

			sortExpr = fmt.Sprintf("ORDER BY %s", pageRequest.Sorting)
		}

		if pageRequest.SkipCount != 0 || pageRequest.MaxResultCount != 0 {
			paginateExpr = "LIMIT @limit OFFSET @offset"
			args["limit"] = pageRequest.MaxResultCount
			args["offset"] = pageRequest.SkipCount
		}

		if err := r.db.QueryRow(ctx, fmt.Sprintf(query, "Count(*)", whereExpr, "", ""), args).Scan(&count); err != nil {
			return nil, 0, fmt.Errorf("unable to count roles: %w", err)
		}
	}

	query = fmt.Sprintf(query, "*", whereExpr, sortExpr, paginateExpr)

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to query roles: %w", err)
	}
	defer rows.Close()

	roles, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Role])

	if count == 0 {
		count = len(roles)
	}

	return roles, count, err
}

func (r roleManager) GetRolesCount(ctx context.Context) (int, error) {
	query := `SELECT Count(*) FROM roles WHERE is_deleted = false`

	var count int

	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to insert row: %w", err)
	}

	return count, nil
}

func (r roleManager) GetRoleById(ctx context.Context, roleId int64) (*models.Role, error) {
	query := `SELECT * FROM roles WHERE id = @roleId AND is_deleted = false LIMIT 1`

	args := pgx.NamedArgs{
		"roleId": roleId,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query role by id: %w", err)
	}
	defer rows.Close()

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Role])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &role, nil
}

func (r roleManager) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	query := `SELECT * FROM roles WHERE normalized_name = @name and is_deleted = false LIMIT 1`

	args := pgx.NamedArgs{
		"name": strings.ToUpper(name),
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query role by name: %w", err)
	}
	defer rows.Close()

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Role])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &role, nil
}

func (r roleManager) GetRolePermission(ctx context.Context, roleId int64, permission string) (*userModel.UserRolePermission, error) {
	query := "SELECT * FROM user_role_permissions WHERE role_id = $1 AND name = $2 LIMIT 1"

	rows, err := r.db.Query(ctx, query, roleId, permission)
	if err != nil {
		return nil, fmt.Errorf("unable to query role permission: %w", err)
	}
	defer rows.Close()

	rolePermission, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userModel.UserRolePermission])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &rolePermission, nil
}

func (r roleManager) GetRolePermissions(ctx context.Context, roleId int64, isGranted *bool) ([]userModel.UserRolePermission, error) {
	query := `
		SELECT urp.* 
		FROM roles r
		JOIN user_role_permissions urp on r.id = urp.role_id
		WHERE r.is_deleted = false AND r.id = @roleId AND (1 = @isGrantedAll OR urp.is_granted = @isGranted) LIMIT 1
	`

	args := pgx.NamedArgs{
		"roleId": roleId,
	}

	if isGranted == nil {
		args["isGrantedAll"] = 1
		args["isGranted"] = true
	} else {
		args["isGrantedAll"] = 0
		args["isGranted"] = *isGranted
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user role permissions: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[userModel.UserRolePermission])
}
