package v1

import (
	"context"
	"net/http"
	"slices"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

// Request
type UpdateUserPermissionsRequest struct {
	Id   int64 `path:"id"`
	Body struct {
		GrantedPermissions []string `json:"grantedPermissions"`
	}
}

// Validator
func (e UpdateUserPermissionsRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Gte(int64(0)).Msg("Invalid user id"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	userRolePermissionManager services.IUserRolePermissionManager,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "UpdateUserPermissions",
			Method:        http.MethodPut,
			Path:          "/user/{id}/permissions",
			Summary:       "Update User Permissions",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationUsersChangePermissions),
				postgres.SetupTransaction(api, pool),
			},
		},
		updateUserPermissions(userRolePermissionManager),
	)
}

func updateUserPermissions(userRolePermissionManager services.IUserRolePermissionManager) func(context.Context, *UpdateUserPermissionsRequest) (*struct{}, error) {
	return func(ctx context.Context, request *UpdateUserPermissionsRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		if err := permissions.ValidatePermissionName(request.Body.GrantedPermissions); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		userManager := services.NewUserManager(tx)

		user, err := userManager.GetUserById(ctx, request.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if user == nil {
			return nil, huma.Error404NotFound("user not found")
		}

		// Get User Role Ids
		roleIds, err := userManager.GetUserRoleIds(ctx, user.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		// Get User Roles Permissions
		userRolesPermissions, err := getUserRolesPermissions(ctx, roleIds, userRolePermissionManager)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		newPermissions := request.Body.GrantedPermissions
		oldPermissions, err := userRolePermissionManager.SetUserPermissions(ctx, user.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		// Prohibit
		for oldPermission := range oldPermissions {
			if !slices.Contains(newPermissions, oldPermission) {
				if err := userManager.DeleteUserPermission(ctx, user.Id, oldPermission); err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}

				if _, ok := userRolesPermissions[oldPermission]; !ok {
					continue
				}

				if err := userManager.CreateUserPermission(ctx, user.Id, oldPermission, false); err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}
			}
		}

		// Granted
		for _, newPermission := range newPermissions {
			if _, ok := oldPermissions[newPermission]; !ok {
				if err := userManager.DeleteUserPermission(ctx, user.Id, newPermission); err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}

				if _, ok := userRolesPermissions[newPermission]; ok {
					continue
				}

				if err := userManager.CreateUserPermission(ctx, user.Id, newPermission, true); err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}
			}
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if _, err := userRolePermissionManager.SetUserPermissions(ctx, user.Id); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	}
}

func getUserRolesPermissions(ctx context.Context, roleIds []int64, userRolePermissionManager services.IUserRolePermissionManager) (map[string]struct{}, error) {
	grantedPermissions := make(map[string]struct{})

	for _, roleId := range roleIds {
		rolePermissions, err := userRolePermissionManager.SetRolePermissions(ctx, roleId)
		if err != nil {
			return nil, err
		}

		for rolePermission := range rolePermissions {
			if _, ok := grantedPermissions[rolePermission]; !ok {
				// add if key does not exists yet
				grantedPermissions[rolePermission] = struct{}{}
			}
		}
	}

	return grantedPermissions, nil
}
