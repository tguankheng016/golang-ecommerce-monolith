package permissions

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"permission_fx",
		permissionProviders,
	)

	permissionProviders = fx.Options(
		fx.Provide(
			NewPermissionManager,
		),
	)
)
