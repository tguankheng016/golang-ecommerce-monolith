package users

import (
	"github.com/tguankheng016/commerce-mono/internal/users/services"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"user_fx",
		userProviders,
	)

	userProviders = fx.Options(
		fx.Provide(
			services.NewUserRolePermissionManager,
		),
	)
)
