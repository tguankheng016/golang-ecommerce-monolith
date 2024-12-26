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
	"github.com/tguankheng016/commerce-mono/pkg/core/pagination"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

// Request
type GetRolesRequest struct {
	pagination.PageRequest
}

// Result
type GetRolesResult struct {
	Body struct {
		pagination.PageResultDto[dtos.RoleDto]
	}
}

// Validator
func (e GetRolesRequest) Schema() v.Schema {
	return v.Schema{
		v.F("skip_count", e.SkipCount):            v.Gte(0).Msg("Page should at least greater than or equal to 0."),
		v.F("max_result_count", e.MaxResultCount): v.Gte(0).Msg("Page size should at least greater than or equal to 0."),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "GetRoles",
			Method:        http.MethodGet,
			Path:          "/roles",
			Summary:       "Get Roles",
			Tags:          []string{"Roles"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationRoles),
			},
		},
		getRoles(pool),
	)
}

func getRoles(pool *pgxpool.Pool) func(context.Context, *GetRolesRequest) (*GetRolesResult, error) {
	return func(ctx context.Context, request *GetRolesRequest) (*GetRolesResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		roleManager := services.NewRoleManager(pool)

		roles, count, err := roleManager.GetRoles(ctx, &request.PageRequest)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var roleDtos []dtos.RoleDto
		if err := copier.Copy(&roleDtos, &roles); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := GetRolesResult{}
		result.Body.PageResultDto = pagination.NewPageResultDto(roleDtos, count)

		return &result, nil
	}
}
