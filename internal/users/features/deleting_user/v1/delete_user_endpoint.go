package v1

import (
	"context"
	"net/http"
	"strings"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/jackc/pgx/v5/pgxpool"
	userConsts "github.com/tguankheng016/commerce-mono/internal/users/constants"
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	httpServer "github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
	"go.uber.org/zap"
)

// Request
type DeleteUserRequest struct {
	Id int64 `path:"id"`
}

// Validator
func (e DeleteUserRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Gte(int64(0)).Msg("Invalid user id"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	cacheManager *cache.Cache[string],
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "DeleteUser",
			Method:        http.MethodDelete,
			Path:          "/user/{id}",
			Summary:       "Delete User",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationUsersDelete),
				postgres.SetupTransaction(api, pool),
			},
		},
		deleteUser(cacheManager),
	)
}

func deleteUser(cacheManager *cache.Cache[string]) func(context.Context, *DeleteUserRequest) (*struct{}, error) {
	return func(ctx context.Context, request *DeleteUserRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		currentUserId, ok := httpServer.GetCurrentUser(ctx)
		if !ok {
			return nil, huma.Error500InternalServerError("current user session not found")
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		userManager := services.NewUserManager(tx)

		user, err := userManager.GetUserById(ctx, request.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if user == nil {
			return nil, huma.Error404NotFound("user not found")
		}

		if user.NormalizedUserName == strings.ToUpper(userConsts.DefaultAdminUserName) {
			return nil, huma.Error400BadRequest("You cannot delete admin's account!")
		}

		if user.Id == currentUserId {
			return nil, huma.Error400BadRequest("You cannot delete your own account!")
		}

		if err := userManager.DeleteUser(ctx, user.Id); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := cacheManager.Delete(ctx, jwt.GenerateStampCacheKey(user.Id)); err != nil {
			// Dont return just log
			logging.Logger.Error("error in deleting security cached", zap.Error(err))
		}

		return nil, nil
	}
}
