package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/core/helpers"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

// Request
type GetUserPermissionsRequest struct {
	Id int64 `path:"id"`
}

// Result
type GetUserPermissionsResult struct {
	Body struct {
		Items []string
	}
}

// Validator
func (e GetUserPermissionsRequest) Schema() v.Schema {
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
			OperationID:   "GetUserPermissons",
			Method:        http.MethodGet,
			Path:          "/user/{id}/permissions",
			Summary:       "Get User Permissions",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationUsersChangePermissions),
			},
		},
		getUserPermissions(userRolePermissionManager),
	)
}

func getUserPermissions(userRolePermissionManager services.IUserRolePermissionManager) func(context.Context, *GetUserPermissionsRequest) (*GetUserPermissionsResult, error) {
	return func(ctx context.Context, request *GetUserPermissionsRequest) (*GetUserPermissionsResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		grantedPermissions, err := userRolePermissionManager.SetUserPermissions(ctx, request.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := GetUserPermissionsResult{}
		result.Body.Items = helpers.MapKeysToSlice(grantedPermissions)

		return &result, nil
	}
}
