package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/tguankheng016/commerce-mono/internal/users/models"
	"github.com/tguankheng016/commerce-mono/pkg/core/pagination"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

type IUserManager interface {
	GetUsers(ctx context.Context, pageRequest *pagination.PageRequest) ([]models.User, int, error)
	GetUsersCount(ctx context.Context) (int, error)
	GetUserById(ctx context.Context, userId int64) (*models.User, error)
	GetUserByUserName(ctx context.Context, userName string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserRoleIds(ctx context.Context, userId int64) ([]int64, error)
	GetUsersInRole(ctx context.Context, roleId int64) ([]models.User, error)
	GetUserPermission(ctx context.Context, userId int64, permission string) (*models.UserRolePermission, error)
	GetUserPermissions(ctx context.Context, userId int64, isGranted *bool) ([]models.UserRolePermission, error)

	CreateUser(ctx context.Context, user *models.User, password string) error
	CreateUserRole(ctx context.Context, userId int64, roleId int64) error
	CreateUserPermission(ctx context.Context, userId int64, permission string, isGranted bool) error

	UpdateUser(ctx context.Context, user *models.User, password string) error
	UpdateUserRoles(ctx context.Context, user *models.User, roles []int64) (bool, error)

	DeleteUser(ctx context.Context, userId int64) error
	DeleteUserRole(ctx context.Context, userId int64, roleId int64) error
	DeleteUserPermission(ctx context.Context, userId int64, permission string) error
	DeleteUserPermissions(ctx context.Context, userId int64) error
}

type userManager struct {
	db postgres.IPgxDbConn
}

func NewUserManager(db postgres.IPgxDbConn) IUserManager {
	return userManager{
		db: db,
	}
}

func (u userManager) GetUsers(ctx context.Context, pageRequest *pagination.PageRequest) ([]models.User, int, error) {
	query := `SELECT %s FROM users WHERE is_deleted = false %s %s %s`
	whereExpr := ""
	sortExpr := ""
	paginateExpr := ""
	count := 0

	args := pgx.NamedArgs{}

	if pageRequest != nil {
		if pageRequest.Filters != "" {
			whereExpr = `
				AND (
					normalized_user_name like @filters OR
					normalized_email like @filters OR
					UPPER(first_name) like @filters OR
					UPPER(last_name) like @filters
				)
			`

			args["filters"] = fmt.Sprintf("%%%s%%", strings.ToUpper(pageRequest.Filters))
		}

		if pageRequest.Sorting != "" {
			sortingFields := []string{"first_name", "last_name", "user_name", "email"}
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

		if err := u.db.QueryRow(ctx, fmt.Sprintf(query, "Count(*)", whereExpr, "", ""), args).Scan(&count); err != nil {
			return nil, 0, fmt.Errorf("unable to count users: %w", err)
		}
	}

	query = fmt.Sprintf(query, "*", whereExpr, sortExpr, paginateExpr)

	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])

	if count == 0 {
		count = len(users)
	}

	return users, count, err
}

func (r userManager) GetUsersCount(ctx context.Context) (int, error) {
	query := `SELECT Count(*) FROM users where is_deleted = false`

	var count int

	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to count row: %w", err)
	}

	return count, nil
}

func (u userManager) GetUserById(ctx context.Context, userId int64) (*models.User, error) {
	query := `SELECT * FROM users where id = @userId and is_deleted = false LIMIT 1`

	args := pgx.NamedArgs{
		"userId": userId,
	}
	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user by id: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (u userManager) GetUserByUserName(ctx context.Context, userName string) (*models.User, error) {
	query := `SELECT * FROM users where normalized_user_name = @userName and is_deleted = false LIMIT 1`

	args := pgx.NamedArgs{
		"userName": strings.ToUpper(userName),
	}
	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user by username: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (u userManager) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT * FROM users where normalized_email = @email and is_deleted = false LIMIT 1`

	args := pgx.NamedArgs{
		"email": strings.ToUpper(email),
	}
	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user by email: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (u userManager) GetUserRoleIds(ctx context.Context, userId int64) ([]int64, error) {
	query := `
		select user_roles.role_id 
		from users join user_roles on users.id = user_roles.user_id 
		where users.is_deleted = false and users.id = @userId
	`

	args := pgx.NamedArgs{
		"userId": userId,
	}

	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user roles id: %w", err)
	}
	defer rows.Close()

	roleIds := make([]int64, 0)

	for rows.Next() {
		var roleId int64
		err := rows.Scan(&roleId)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		roleIds = append(roleIds, roleId)
	}

	return roleIds, nil
}

func (u userManager) GetUsersInRole(ctx context.Context, roleId int64) ([]models.User, error) {
	query := "SELECT DISTINCT B.* FROM user_roles A JOIN users B on A.user_id = B.id WHERE A.role_id = $1"

	rows, err := u.db.Query(ctx, query, roleId)
	if err != nil {
		return nil, fmt.Errorf("unable to query users in role: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
}

func (u userManager) GetUserPermission(ctx context.Context, userId int64, permission string) (*models.UserRolePermission, error) {
	query := "SELECT * FROM user_role_permissions WHERE user_id = $1 AND name = $2 LIMIT 1"

	rows, err := u.db.Query(ctx, query, userId, permission)
	if err != nil {
		return nil, fmt.Errorf("unable to query user permission: %w", err)
	}
	defer rows.Close()

	userPermission, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.UserRolePermission])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &userPermission, nil
}

func (u userManager) GetUserPermissions(ctx context.Context, userId int64, isGranted *bool) ([]models.UserRolePermission, error) {
	query := `
		SELECT urp.* 
		FROM users u
		JOIN user_role_permissions urp on u.id = urp.user_id
		WHERE u.is_deleted = false AND u.id = @userId AND (1 = @isGrantedAll OR urp.is_granted = @isGranted) LIMIT 1
	`

	args := pgx.NamedArgs{
		"userId": userId,
	}

	if isGranted == nil {
		args["isGrantedAll"] = 1
		args["isGranted"] = true
	} else {
		args["isGrantedAll"] = 0
		args["isGranted"] = *isGranted
	}

	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user permissions: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.UserRolePermission])
}
