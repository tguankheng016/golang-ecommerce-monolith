package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/commerce-mono/internal/identities/dtos"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	httpServer "github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

// Result
type GetCurrentSessionResult struct {
	User               *dtos.UserLoginInfoDto `json:"user"`
	AllPermissions     map[string]bool        `json:"allPermissions"`
	GrantedPermissions map[string]bool        `json:"grantedPermissions"`
}

type GetHumaCurrentSessionResult struct {
	Body struct {
		GetCurrentSessionResult
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	userPermissionManager userService.IUserRolePermissionManager,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "GetCurrentSession",
			Method:        http.MethodGet,
			Path:          "/identities/current-session",
			Summary:       "Get Current Session",
			Tags:          []string{"Identities"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
		},
		getAllPermissions(pool, userPermissionManager),
	)
}

func getAllPermissions(pool *pgxpool.Pool, userPermissionManager userService.IUserRolePermissionManager) func(context.Context, *struct{}) (*GetHumaCurrentSessionResult, error) {
	return func(ctx context.Context, input *struct{}) (*GetHumaCurrentSessionResult, error) {
		userId, ok := httpServer.GetCurrentUser(ctx)
		result := GetHumaCurrentSessionResult{}

		allPermisions := permissions.GetAppPermissions().Items
		result.Body.AllPermissions = convertPermissionMap(allPermisions)

		if ok {
			userManager := userService.NewUserManager(pool)
			user, err := userManager.GetUserById(ctx, userId)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
			if user == nil {
				return nil, huma.Error404NotFound("user not found")
			}

			currentUserLoginInfoDto := dtos.UserLoginInfoDto{}
			if err := copier.Copy(&currentUserLoginInfoDto, &user); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			result.Body.User = &currentUserLoginInfoDto

			userGrantedPermissions, err := userPermissionManager.SetUserPermissions(ctx, userId)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			result.Body.GrantedPermissions = convertPermissionMap(userGrantedPermissions)
		} else {
			result.Body.GrantedPermissions = make(map[string]bool)
		}

		return &result, nil
	}
}

func convertPermissionMap[V any](permissionMap map[string]V) map[string]bool {
	boolMap := make(map[string]bool)
	for key := range permissionMap {
		boolMap[key] = true
	}
	return boolMap
}
