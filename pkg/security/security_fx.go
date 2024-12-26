package security

import (
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"security_fx",
		securityProviders,
	)

	securityProviders = fx.Options(
		fx.Provide(
			jwt.NewTokenHandler,
			jwt.NewTokenKeyValidator,
			jwt.NewSecurityStampValidator,
		),
	)
)
