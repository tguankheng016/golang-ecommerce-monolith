package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/commerce-mono/internal/users/dtos"
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/core/pagination"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

// Request
type GetUsersRequest struct {
	pagination.PageRequest
}

// Result
type GetUsersResult struct {
	Body struct {
		pagination.PageResultDto[dtos.UserDto]
	}
}

// Validator
func (e GetUsersRequest) Schema() v.Schema {
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
			OperationID:   "GetUsers",
			Method:        http.MethodGet,
			Path:          "/users",
			Summary:       "Get Users",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationUsers),
			},
		},
		getUsers(pool),
	)
}

func getUsers(pool *pgxpool.Pool) func(context.Context, *GetUsersRequest) (*GetUsersResult, error) {
	return func(ctx context.Context, request *GetUsersRequest) (*GetUsersResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		userManager := services.NewUserManager(pool)

		users, count, err := userManager.GetUsers(ctx, &request.PageRequest)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var userDtos []dtos.UserDto
		if err := copier.Copy(&userDtos, &users); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := GetUsersResult{}
		result.Body.PageResultDto = pagination.NewPageResultDto(userDtos, count)

		return &result, nil
	}
}
