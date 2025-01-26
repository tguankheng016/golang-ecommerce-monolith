package v1

import (
	"context"
	"net/http"
	"regexp"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	roleService "github.com/tguankheng016/commerce-mono/internal/roles/services"
	"github.com/tguankheng016/commerce-mono/internal/users/dtos"
	"github.com/tguankheng016/commerce-mono/internal/users/models"
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

// Request
type HumaCreateUserRequest struct {
	Body struct {
		dtos.CreateUserDto
	}
}

// Result
type HumaCreateUserResult struct {
	Body struct {
		User dtos.UserDto
	}
}

// Validator
func (e HumaCreateUserRequest) Schema() v.Schema {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return v.Schema{
		v.F("id", e.Body.Id): v.Any(
			v.Zero[*int64](),
			v.Nested(func(ptr *int64) v.Validator { return v.Value(*ptr, v.Eq(int64(0)).Msg("Invalid user id")) }),
		).LastError(),
		v.F("username", e.Body.UserName):    v.Nonzero[string]().Msg("Please enter the username"),
		v.F("first_name", e.Body.FirstName): v.Nonzero[string]().Msg("Please enter the first name"),
		v.F("last_name", e.Body.LastName):   v.Nonzero[string]().Msg("Please enter the last name"),
		v.F("email", e.Body.Email): v.All(
			v.Nonzero[string]().Msg("Please enter the email"),
			v.Match(pattern).Msg("Please enter a valid email"),
		),
		v.F("password", e.Body.Password): v.Nonzero[string]().Msg("Please enter the password"),
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
			OperationID:   "CreateUser",
			Method:        http.MethodPost,
			Path:          "/user",
			Summary:       "Create User",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationUsersCreate),
				postgres.SetupTransaction(api, pool),
			},
		},
		createUser(),
	)
}

func createUser() func(context.Context, *HumaCreateUserRequest) (*HumaCreateUserResult, error) {
	return func(ctx context.Context, request *HumaCreateUserRequest) (*HumaCreateUserResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		userManager := services.NewUserManager(tx)

		var user models.User
		if err := copier.Copy(&user, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := userManager.CreateUser(ctx, &user, request.Body.Password); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		if len(request.Body.RoleIds) > 0 {
			roleManager := roleService.NewRoleManager(tx)

			for _, roleId := range request.Body.RoleIds {
				if roleId == 0 {
					continue
				}

				role, err := roleManager.GetRoleById(ctx, roleId)
				if err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}
				if role == nil {
					return nil, huma.Error400BadRequest("Invalid role id")
				}

				if err := userManager.CreateUserRole(ctx, user.Id, roleId); err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}
			}
		}

		var userDto dtos.UserDto
		if err := copier.Copy(&userDto, &user); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := HumaCreateUserResult{}
		result.Body.User = userDto

		return &result, nil
	}
}
