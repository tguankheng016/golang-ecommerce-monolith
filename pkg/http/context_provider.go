package http

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tguankheng016/commerce-mono/pkg/logging"
)

// NewContext returns a context that listens for SIGINT and SIGTERM signals.
// When the signals are received, the context is canceled and the logger logs
// a message.
func NewContext() context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			<-ctx.Done()
			logging.Logger.Info("context is canceled!")
			cancel()
			return
		}
	}()

	return ctx
}
