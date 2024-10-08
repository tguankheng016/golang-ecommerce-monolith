package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"gorm.io/gorm"
)

func TransactionalContextMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tx := db.Begin()
			if tx.Error != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, tx.Error.Error())
			}

			ctx := context.WithValue(c.Request().Context(), constants.TxKey(constants.DbContextKey), tx)
			c.SetRequest(c.Request().WithContext(ctx))

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

			if tx.Commit().Error != nil {
				log.Println("Failed to commit transaction")
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
			}

			return nil
		}
	}
}
