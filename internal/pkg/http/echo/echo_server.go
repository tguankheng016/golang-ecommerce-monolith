package echoserver

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

func NewEchoServer() *echo.Echo {
	e := echo.New()
	return e
}

func RunHttpServer(ctx context.Context, echo *echo.Echo, log logger.ILogger, cfg *EchoOptions) error {
	echo.Server.ReadTimeout = ReadTimeout
	echo.Server.WriteTimeout = WriteTimeout
	echo.Server.MaxHeaderBytes = MaxHeaderBytes

	go func() {
		<-ctx.Done()
		log.Infof("shutting down Http PORT: {%s}", cfg.Port)
		err := echo.Shutdown(ctx)
		if err != nil {
			log.Errorf("(Shutdown) err: {%v}", err)
			return
		}
		log.Info("server exited properly")
	}()

	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			log.Infof("shutting down Http PORT: {%s}", cfg.Port)
	// 			err := echo.Shutdown(ctx)
	// 			if err != nil {
	// 				log.Errorf("(Shutdown) err: {%v}", err)
	// 				return
	// 			}
	// 			log.Info("server exited properly")
	// 			return
	// 		}
	// 	}
	// }()

	err := echo.Start(cfg.Port)

	return err
}
