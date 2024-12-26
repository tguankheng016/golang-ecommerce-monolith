package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/commerce-mono/internal/roles/dtos"
	"github.com/tguankheng016/commerce-mono/internal/roles/models"
	"github.com/tguankheng016/commerce-mono/internal/roles/services"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

// Request
type HumaCreateRoleRequest struct {
	Body struct {
		dtos.CreateRoleDto
	}
}

// Result
type HumaCreateRoleResult struct {
	Body struct {
		Role dtos.RoleDto
	}
}

// Validator
func (e HumaCreateRoleRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Body.Id): v.Any(
			v.Zero[*int64](),
			v.Nested(func(ptr *int64) v.Validator { return v.Value(*ptr, v.Eq(int64(0)).Msg("Invalid role id")) }),
		).LastError(),
		v.F("rolename", e.Body.Name): v.Nonzero[string]().Msg("Please enter the role name"),
	}
}

// Handler
// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "CreateRole",
			Method:        http.MethodPost,
			Path:          "/role",
			Summary:       "Create Role",
			Tags:          []string{"Roles"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationRolesCreate),
				postgres.SetupTransaction(api, pool),
			},
		},
		createRole(),
	)
}

func createRole() func(context.Context, *HumaCreateRoleRequest) (*HumaCreateRoleResult, error) {
	return func(ctx context.Context, request *HumaCreateRoleRequest) (*HumaCreateRoleResult, error) {
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

		roleManager := services.NewRoleManager(tx)

		var role models.Role
		if err := copier.Copy(&role, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := roleManager.CreateRole(ctx, &role); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		if len(request.Body.GrantedPermissions) > 0 {
			for _, grantedPermission := range request.Body.GrantedPermissions {
				if err := roleManager.CreateRolePermission(ctx, role.Id, grantedPermission, true); err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}
			}
		}

		var roleDto dtos.RoleDto
		if err := copier.Copy(&roleDto, &role); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := HumaCreateRoleResult{}
		result.Body.Role = roleDto

		return &result, nil
	}
}
