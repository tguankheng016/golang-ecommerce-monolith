package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/config"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"go.uber.org/fx"
)

func RunServers(lc fx.Lifecycle, log logger.ILogger, e *echo.Echo, ctx context.Context, cfg *config.Config) error {

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := echoServer.RunHttpServer(ctx, e, log, cfg.EchoOptions); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running http server: %v", err)
				}
			}()

			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "Golang ECommerce")
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Infof("all servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
