package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/commerce-mono/internal/roles/dtos"
	"github.com/tguankheng016/commerce-mono/internal/roles/services"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/core/helpers"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

// Request
type GetRoleByIdRequest struct {
	Id int64 `path:"id"`
}

// Result
type GetRoleByIdResult struct {
	Body struct {
		Role dtos.CreateOrEditRoleDto
	}
}

// Validator
func (e GetRoleByIdRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Gte(int64(0)).Msg("Invalid role id"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	permissionManager userService.IUserRolePermissionManager,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "GetRoleById",
			Method:        http.MethodGet,
			Path:          "/role/{id}",
			Summary:       "Get Role By Id",
			Tags:          []string{"Roles"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationRoles),
			},
		},
		getRoleById(pool, permissionManager),
	)
}

func getRoleById(pool *pgxpool.Pool, permissionManager userService.IUserRolePermissionManager) func(context.Context, *GetRoleByIdRequest) (*GetRoleByIdResult, error) {
	return func(ctx context.Context, request *GetRoleByIdRequest) (*GetRoleByIdResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		var roleEditDto dtos.CreateOrEditRoleDto

		if request.Id > 0 {
			roleManager := services.NewRoleManager(pool)

			role, err := roleManager.GetRoleById(ctx, request.Id)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
			if role == nil {
				return nil, huma.Error404NotFound("role not found")
			}
			if err := copier.Copy(&roleEditDto, &role); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			rolePermissions, err := permissionManager.SetRolePermissions(ctx, role.Id)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			roleEditDto.GrantedPermissions = helpers.MapKeysToSlice(rolePermissions)
		} else {
			roleEditDto = dtos.CreateOrEditRoleDto{}
		}

		result := GetRoleByIdResult{}
		result.Body.Role = roleEditDto

		return &result, nil
	}
}
