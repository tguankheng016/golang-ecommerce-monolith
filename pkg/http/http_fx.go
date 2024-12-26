package http

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"http_fx",
		httpProviders,
		httpInvokes,
	)

	httpProviders = fx.Options(
		fx.Provide(
			NewContext,
			NewHumaRouter,
			NewHumaServer,
			NewHumaListener,
		),
	)

	httpInvokes = fx.Options(
		fx.Invoke(RunHumaServers),
	)
)
