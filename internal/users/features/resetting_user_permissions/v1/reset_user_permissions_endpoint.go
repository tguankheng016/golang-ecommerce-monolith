package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

// Request
type ResetUserPermissionsRequest struct {
	Id int64 `path:"id"`
}

// Validator
func (e ResetUserPermissionsRequest) Schema() v.Schema {
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
			OperationID:   "ResetUserPermissions",
			Method:        http.MethodPut,
			Path:          "/user/{id}/reset-permissions",
			Summary:       "Reset User Permissions",
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
		resetUserPermissions(userRolePermissionManager),
	)
}

func resetUserPermissions(userRolePermissionManager services.IUserRolePermissionManager) func(context.Context, *ResetUserPermissionsRequest) (*struct{}, error) {
	return func(ctx context.Context, request *ResetUserPermissionsRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
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

		if err := userManager.DeleteUserPermissions(ctx, request.Id); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
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
