package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/tguankheng016/commerce-mono/pkg/environment"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHumaRouter() *chi.Mux {
	router := chi.NewMux()
	return router
}

func NewHumaListener(env environment.Environment, options *ServerOptions) (net.Listener, error) {
	if !env.IsTest() {
		return nil, nil
	}

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	address := ln.Addr().String()

	parts := strings.Split(address, ":")
	port := parts[len(parts)-1]
	options.Port = port

	return ln, nil
}

func NewHumaServer(router *chi.Mux, options *ServerOptions) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", options.Host, options.Port),
		Handler: router,
	}

	return server
}

func RunHumaServer(env environment.Environment, server *http.Server, cfg *ServerOptions, ln net.Listener) error {
	if env.IsTest() {
		logging.Logger.Info("huma server started ...", zap.String("baseUrl", cfg.GetBasePath()))
		err := server.Serve(ln)
		return err
	} else {
		logging.Logger.Info("huma server started ...", zap.String("baseUrl", cfg.GetBasePath()))
		err := server.ListenAndServe()
		return err
	}
}

func RunHumaServers(lc fx.Lifecycle, server *http.Server, cfg *ServerOptions, env environment.Environment, ln net.Listener) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := RunHumaServer(env, server, cfg, ln); !errors.Is(err, http.ErrServerClosed) {
					logging.Logger.Fatal("error running http server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logging.Logger.Info("shutting down Http PORT: " + cfg.Port)

			if err := server.Shutdown(ctx); err != nil {
				logging.Logger.Error("(Shutdown) err: ", zap.Error(err))
			}

			logging.Logger.Info("all http servers shutdown gracefully...")

			return nil
		},
	})

	return nil
}
