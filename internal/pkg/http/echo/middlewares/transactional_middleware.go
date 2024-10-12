package middlewares

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"gorm.io/gorm"
)

// SetupTransaction returns an Echo middleware that sets up a transactional context for the incoming request.
// The transaction is started at the beginning of the request and rolled back if any error is returned from the handler.
// If no error is returned, the transaction is committed at the end of the request.
// The transaction is also rolled back if a panic is detected.
// The transaction is stored in the context under the key constants.DbContextKey.
// The middleware uses the provided skipper to determine if the request should be skipped.
// If the skipper returns true, the middleware does not start a transaction and instead calls the next handler in the chain.
func SetupTransaction(skipper echoMiddleware.Skipper, db *gorm.DB) echo.MiddlewareFunc {
	// Defaults
	if skipper == nil {
		skipper = echoMiddleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			ctx := c.Request().Context()
			userId, ok := echoServer.GetCurrentUser(c)
			if ok {
				ctx = context.WithValue(ctx, constants.CtxKey(constants.CurrentUserContextKey), userId)
			}

			tx := db.WithContext(ctx).Begin()
			if tx.Error != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, tx.Error.Error())
			}

			c.Set(constants.DbContextKey, tx)
			//ctx := context.WithValue(c.Request().Context(), constants.CtxKey(constants.DbContextKey), tx)
			//c.SetRequest(c.Request().WithContext(ctx))

			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
					log.Println("Recovered from panic, transaction rolled back")
				}
			}()

			err := next(c)
			if err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Commit().Error; err != nil && err != sql.ErrTxDone {
				log.Println("Failed to commit transaction")
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
			}

			return nil
		}
	}
}
