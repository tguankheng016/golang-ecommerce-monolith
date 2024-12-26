package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/roles/services"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

// Request
type DeleteRoleRequest struct {
	Id int64 `path:"id"`
}

// Validator
func (e DeleteRoleRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Gte(int64(0)).Msg("Invalid role id"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	userRolePermissionManager userService.IUserRolePermissionManager,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "DeleteRole",
			Method:        http.MethodDelete,
			Path:          "/role/{id}",
			Summary:       "Delete Role",
			Tags:          []string{"Roles"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationRolesDelete),
				postgres.SetupTransaction(api, pool),
			},
		},
		deleteRole(userRolePermissionManager),
	)
}

func deleteRole(userRolePermissionManager userService.IUserRolePermissionManager) func(context.Context, *DeleteRoleRequest) (*struct{}, error) {
	return func(ctx context.Context, request *DeleteRoleRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		roleManager := services.NewRoleManager(tx)
		userManager := userService.NewUserManager(tx)

		role, err := roleManager.GetRoleById(ctx, request.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if role == nil {
			return nil, huma.Error404NotFound("role not found")
		}

		if role.IsStatic {
			return nil, huma.Error400BadRequest("You cannot delete static role!")
		}

		users, err := userManager.GetUsersInRole(ctx, role.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		for _, user := range users {
			if err := userManager.DeleteUserRole(ctx, user.Id, role.Id); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			userRolePermissionManager.RemoveUserRoleCaches(ctx, user.Id)
		}

		if err := roleManager.DeleteRole(ctx, role.Id); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	}
}
