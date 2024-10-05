package configurations

import (
	"context"

	"github.com/labstack/echo/v4"
	getting_users "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/features/getting_users/v1/endpoints"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"gorm.io/gorm"
)

func ConfigEndpoints(db *gorm.DB, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	getting_users.MapRoute(db, log, echo, ctx)
}
