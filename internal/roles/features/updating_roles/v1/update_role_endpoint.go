package v1

import (
	"context"
	"net/http"
	"slices"
	"strings"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	roleConsts "github.com/tguankheng016/commerce-mono/internal/roles/constants"
	"github.com/tguankheng016/commerce-mono/internal/roles/dtos"
	"github.com/tguankheng016/commerce-mono/internal/roles/services"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

// Request
type HumaUpdateRoleRequest struct {
	Body struct {
		dtos.EditRoleDto
	}
}

// Result
type HumaUpdateRoleResult struct {
	Body struct {
		Role dtos.RoleDto
	}
}

// Validator
func (e HumaUpdateRoleRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Body.Id): v.All(
			v.Nonzero[*int64]().Msg("Invalid role id"),
			v.Nested(func(ptr *int64) v.Validator { return v.Value(*ptr, v.Gt(int64(0)).Msg("Invalid role id")) }),
		),
		v.F("rolename", e.Body.Name): v.Nonzero[string]().Msg("Please enter the role name"),
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
			OperationID:   "UpdateRole",
			Method:        http.MethodPut,
			Path:          "/role",
			Summary:       "Update Role",
			Tags:          []string{"Roles"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationRolesEdit),
				postgres.SetupTransaction(api, pool),
			},
		},
		updateRole(userRolePermissionManager),
	)
}

func updateRole(userRolePermissionManager userService.IUserRolePermissionManager) func(context.Context, *HumaUpdateRoleRequest) (*HumaUpdateRoleResult, error) {
	return func(ctx context.Context, request *HumaUpdateRoleRequest) (*HumaUpdateRoleResult, error) {
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

		role, err := roleManager.GetRoleById(ctx, *request.Body.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if role == nil {
			return nil, huma.Error404NotFound("role not found")
		}
		if role.IsStatic && role.NormalizedName != strings.ToUpper(request.Body.Name) {
			return nil, huma.Error400BadRequest("You cannot change the name of static role")
		}

		if err := copier.Copy(&role, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := roleManager.UpdateRole(ctx, role); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		oldPermissions, err := userRolePermissionManager.SetRolePermissions(ctx, role.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		isAdmin := strings.EqualFold(role.Name, roleConsts.DefaultAdminRoleName)

		// Prohibit
		for oldPermission := range oldPermissions {
			if !slices.Contains(request.Body.GrantedPermissions, oldPermission) {
				if err := roleManager.DeleteRolePermission(ctx, role.Id, oldPermission); err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}

				if isAdmin {
					if err := roleManager.CreateRolePermission(ctx, role.Id, oldPermission, false); err != nil {
						return nil, huma.Error500InternalServerError(err.Error())
					}
				}
			}
		}

		// Granted
		for _, newPermission := range request.Body.GrantedPermissions {
			if _, ok := oldPermissions[newPermission]; !ok {
				rolePermission, err := roleManager.GetRolePermission(ctx, role.Id, newPermission)
				if err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}

				if rolePermission == nil {
					if err := roleManager.CreateRolePermission(ctx, role.Id, newPermission, true); err != nil {
						return nil, huma.Error500InternalServerError(err.Error())
					}
				} else if isAdmin && !rolePermission.IsGranted {
					if err := roleManager.DeleteRolePermission(ctx, role.Id, newPermission); err != nil {
						return nil, huma.Error500InternalServerError(err.Error())
					}
				}
			}
		}

		var roleDto dtos.RoleDto
		if err := copier.Copy(&roleDto, &role); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := HumaUpdateRoleResult{}
		result.Body.Role = roleDto

		if err := tx.Commit(ctx); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if _, err := userRolePermissionManager.SetRolePermissions(ctx, role.Id); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &result, nil
	}
}
