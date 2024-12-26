package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	httpServer "github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"github.com/tguankheng016/commerce-mono/pkg/security"
)

// Request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
type HumaChangePasswordRequest struct {
	Body struct {
		ChangePasswordRequest
	}
}

// Validator
func (e HumaChangePasswordRequest) Schema() v.Schema {
	return v.Schema{
		v.F("current_password", e.Body.CurrentPassword): v.Nonzero[string]().Msg("Please enter the current password"),
		v.F("new_password", e.Body.NewPassword):         v.Nonzero[string]().Msg("Please enter the new password"),
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
			OperationID:   "ChangePassword",
			Method:        http.MethodPut,
			Path:          "/identities/change-password",
			Summary:       "Change Password",
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
		changePassword(),
	)
}

func changePassword() func(context.Context, *HumaChangePasswordRequest) (*struct{}, error) {
	return func(ctx context.Context, request *HumaChangePasswordRequest) (*struct{}, error) {
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

		ok, err = security.ComparePasswords(user.PasswordHash, request.Body.CurrentPassword)
		if err != nil || !ok {
			return nil, huma.Error400BadRequest("current password is incorrect")
		}

		if err := userManager.UpdateUser(ctx, user, request.Body.NewPassword); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, nil
	}
}
