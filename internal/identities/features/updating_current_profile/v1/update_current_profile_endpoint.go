package v1

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	httpServer "github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
)

// Request
type UpdateCurrentProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
	Email     string `json:"Email"`
}
type HumaUpdateCurrentProfileRequest struct {
	Body struct {
		UpdateCurrentProfileRequest
	}
}

// Validator
func (e HumaUpdateCurrentProfileRequest) Schema() v.Schema {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return v.Schema{
		v.F("firstName", e.Body.FirstName): v.Nonzero[string]().Msg("Please enter your first name"),
		v.F("lastName", e.Body.LastName):   v.Nonzero[string]().Msg("Please enter your last name"),
		v.F("userName", e.Body.UserName):   v.Nonzero[string]().Msg("Please enter your username"),
		v.F("email", e.Body.Email): v.All(
			v.Nonzero[string]().Msg("Please enter your email"),
			v.Match(pattern).Msg("Please enter a valid email"),
		),
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
			OperationID:   "UpdateCurrentProfile",
			Method:        http.MethodPut,
			Path:          "/identities/current-profile",
			Summary:       "Update Current Profile",
			Tags:          []string{"Identities"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, ""),
				postgres.SetupTransaction(api, pool),
			},
		},
		updateCurrentProfile(),
	)
}

func updateCurrentProfile() func(context.Context, *HumaUpdateCurrentProfileRequest) (*struct{}, error) {
	return func(ctx context.Context, request *HumaUpdateCurrentProfileRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		userId, ok := httpServer.GetCurrentUser(ctx)
		if !ok {
			return nil, huma.Error400BadRequest("current user not found")
		}

		userManager := services.NewUserManager(tx)

		user, err := userManager.GetUserById(ctx, userId)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if user == nil {
			return nil, huma.Error404NotFound("user not found")
		}

		if user.NormalizedUserName == strings.ToUpper(userConsts.DefaultAdminUserName) && user.UserName != request.Body.UserName {
			return nil, huma.Error400BadRequest("You cannot update admin's username!")
		}

		if err := copier.Copy(&user, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := userManager.UpdateUser(ctx, user, ""); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, nil
	}
}
