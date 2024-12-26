package configurations

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"configuration_fx",
		configurationInvokes,
	)

	configurationInvokes = fx.Options(
		//fx.Invoke(ConfigMiddlewares),
		// fx.Invoke(ConfigSwagger),
		//fx.Invoke(ConfigEndpoints),
		fx.Invoke(ConfigEndpoints),
	// fx.Invoke(ConfigUserGrpcServer),
	// fx.Invoke(ConfigIdentityGrpcServer),
	// fx.Invoke(ConfigPermissionGrpcServer),
	)
)
