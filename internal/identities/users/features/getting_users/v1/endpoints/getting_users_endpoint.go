package endpoints

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/mapper"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/pagination"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/permissions"
	"gorm.io/gorm"
)

type GetUsersRequestDto struct {
	*pagination.PageRequest
}

type GetUsersResponseDto struct {
	*pagination.PageResultDto[dtos.UserDto]
} // @name GetUsersResponseDto

func MapRoute(db *gorm.DB, jwt jwt.IJwtTokenValidator, checker permissions.IPermissionChecker, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/users")
	group.GET("", getAllUsers(db, log, ctx), middlewares.ValidateToken(jwt), middlewares.Authorize(checker, permissions.PagesAdministrationUsers))
}

// GetAllUsers
// @Tags Users
// @Summary Get all users
// @Description Get all users
// @Accept json
// @Produce json
// @Param GetUsersRequestDto query GetUsersRequestDto false "GetUsersRequestDto"
// @Success 200 {object} GetUsersResponseDto
// @Security ApiKeyAuth
// @Router /api/v1/users [get]
func getAllUsers(db *gorm.DB, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []models.User

		pageRequest, err := pagination.GetPageRequestFromCtx(c)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		fields := []string{"first_name", "last_name", "user_name", "email"}

		if err := pageRequest.SanitizeSorting(fields...); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		userPageRequest := &GetUsersRequestDto{PageRequest: pageRequest}

		query := db
		countQuery := db.Model(&models.User{})

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
			log.Warnf("GetUsers", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		userDtos, err := mapper.Map[[]dtos.UserDto](users)

		if err != nil {
			log.Warnf("BindUsers", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result := &GetUsersResponseDto{
			pagination.NewPageResultDto(userDtos, count),
		}

		return c.JSON(http.StatusOK, result)
	}
}
