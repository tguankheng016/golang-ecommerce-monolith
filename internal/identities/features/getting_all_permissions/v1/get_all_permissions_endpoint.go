package v1

import (
	"context"
	"net/http"
	"sort"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/commerce-mono/internal/identities/dtos"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
)

// Result
type GetAllPermissionsResult struct {
	Body struct {
		Items []dtos.PermissionGroupDto
	}
}

// Handler
func MapRoute(
	api huma.API,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "GetAllPermissions",
			Method:        http.MethodGet,
			Path:          "/identities/permissions",
			Summary:       "Get All Permissions",
			Tags:          []string{"Identities"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, ""),
			},
		},
		getAllPermissions(),
	)
}

func getAllPermissions() func(context.Context, *struct{}) (*GetAllPermissionsResult, error) {
	return func(ctx context.Context, input *struct{}) (*GetAllPermissionsResult, error) {
		allPermissions := permissions.GetAppPermissions().Items

		groupedPermissions := make(map[string][]permissions.Permission)

		for _, permission := range allPermissions {
			groupedPermissions[permission.Group] = append(groupedPermissions[permission.Group], permission)
		}

		var permissionGroups []dtos.PermissionGroupDto
		for groupName, permissions := range groupedPermissions {
			var permissionDtos []dtos.PermissionDto
			if err := copier.Copy(&permissionDtos, &permissions); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			sort.Slice(permissionDtos, func(i, j int) bool {
				return permissionDtos[i].Name < permissionDtos[j].Name
			})

			permissionGroups = append(permissionGroups, dtos.PermissionGroupDto{
				GroupName:   groupName,
				Permissions: permissionDtos,
			})
		}

		sort.Slice(permissionGroups, func(i, j int) bool {
			return permissionGroups[i].GroupName < permissionGroups[j].GroupName
		})

		result := GetAllPermissionsResult{}
		result.Body.Items = permissionGroups

		return &result, nil
	}
}
