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
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

// Request
type GetUserByIdRequest struct {
	Id int64 `path:"id"`
}

// Result
type GetUserByIdResult struct {
	Body struct {
		User dtos.CreateOrEditUserDto
	}
}

// Validator
func (e GetUserByIdRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Gte(int64(0)).Msg("Invalid user id"),
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
			OperationID:   "GetUserById",
			Method:        http.MethodGet,
			Path:          "/user/{id}",
			Summary:       "Get User By Id",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationUsers),
			},
		},
		getUserById(pool),
	)
}

func getUserById(pool *pgxpool.Pool) func(context.Context, *GetUserByIdRequest) (*GetUserByIdResult, error) {
	return func(ctx context.Context, request *GetUserByIdRequest) (*GetUserByIdResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		var userEditDto dtos.CreateOrEditUserDto

		if request.Id > 0 {
			userManager := services.NewUserManager(pool)

			user, err := userManager.GetUserById(ctx, request.Id)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
			if user == nil {
				return nil, huma.Error404NotFound("user not found")
			}
			if err := copier.Copy(&userEditDto, &user); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
			userEditDto.RoleIds, err = userManager.GetUserRoleIds(ctx, user.Id)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

		} else {
			userEditDto = dtos.CreateOrEditUserDto{}
			userEditDto.RoleIds = make([]int64, 0)
		}

		result := GetUserByIdResult{}
		result.Body.User = userEditDto

		return &result, nil
	}
}
