package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/pagination"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
)

type GetUsersRequest struct {
	*pagination.PageRequest
}

type GetUsersResult struct {
	*pagination.PageResultDto[dtos.UserDto]
} // @name GetUsersResult

func MapRoute(echo *echo.Echo, validator *validator.Validate) {
	group := echo.Group("/api/v1/users")
	group.GET("", getAllUsers(validator), middlewares.Authorize(permissions.PagesAdministrationUsers))
}

// GetAllUsers
// @Tags Users
// @Summary Get all users
// @Description Get all users
// @Accept json
// @Produce json
// @Param GetUsersRequest query GetUsersRequest false "GetUsersRequest"
// @Success 200 {object} GetUsersResult
// @Security ApiKeyAuth
// @Router /api/v1/users [get]
func getAllUsers(validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := database.GetTxFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var users []models.User

		pageRequest, err := pagination.GetPageRequestFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, pageRequest); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		fields := []string{"first_name", "last_name", "user_name", "email"}

		if err := pageRequest.SanitizeSorting(fields...); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		userPageRequest := &GetUsersRequest{PageRequest: pageRequest}

		query := tx
		countQuery := tx.Model(&models.User{})

		if userPageRequest.Filters != "" {
			likeExpr := userPageRequest.BuildFiltersExpr(fields...)
			query = query.Where(likeExpr)
			countQuery = countQuery.Where(likeExpr)
		}

		if userPageRequest.Sorting != "" {
			query = query.Order(userPageRequest.Sorting)
		}

		if userPageRequest.SkipCount > 0 || userPageRequest.MaxResultCount > 0 {
			query = userPageRequest.Paginate(query)
		}

		var count int64

		if err := countQuery.Count(&count).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := query.Find(&users).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var userDtos []dtos.UserDto
		if err := copier.Copy(&userDtos, &users); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result := &GetUsersResult{
			pagination.NewPageResultDto(userDtos, count),
		}

		return c.JSON(http.StatusOK, result)
	}
}
