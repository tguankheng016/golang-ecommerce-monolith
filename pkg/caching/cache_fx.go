package caching

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"caching_fx",
		cachingProvider,
		cachingInvoke,
	)

	cachingProvider = fx.Options(
		fx.Provide(
			NewCacheManager,
		),
	)

	cachingInvoke = fx.Options(
		fx.Invoke(RunCaching),
	)
)
