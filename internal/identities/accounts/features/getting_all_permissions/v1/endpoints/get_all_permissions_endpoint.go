package endpoints

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

type PermissionDto struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IsGranted   bool   `json:"isGranted"`
} // @name PermissionDto

type PermissionGroupDto struct {
	GroupName   string          `json:"groupName"`
	Permissions []PermissionDto `json:"permissions"`
} // @name PermissionGroupDto

type GetAllPermissionResult struct {
	Items []PermissionGroupDto `json:"items"`
} // @name GetAllPermissionResult

func MapRoute(echo *echo.Echo) {
	group := echo.Group("/api/v1/accounts/app-permissions")
	group.GET("", getAllPermissions())
}

// Get GetAllAppPermissions
// @Tags Accounts
// @Summary Get All App Permissions
// @Description Get All App Permissions
// @Accept json
// @Produce json
// @Success 200 {object} GetAllPermissionResult
// @Security ApiKeyAuth
// @Router /api/v1/accounts/app-permissions [get]
func getAllPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		allPermisions := permissions.GetAppPermissions().Items

		groupedPermissions := make(map[string][]permissions.Permission)

		for _, permission := range allPermisions {
			groupedPermissions[permission.Group] = append(groupedPermissions[permission.Group], permission)
		}

		var result []PermissionGroupDto
		for groupName, permissions := range groupedPermissions {
			var permissionDtos []PermissionDto
			if err := copier.Copy(&permissionDtos, &permissions); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			result = append(result, PermissionGroupDto{
				GroupName:   groupName,
				Permissions: permissionDtos,
			})
		}

		return c.JSON(http.StatusOK, GetAllPermissionResult{Items: result})
	}
}
