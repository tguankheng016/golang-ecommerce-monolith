package main

import (
	"github.com/tguankheng016/commerce-mono/config"
	"github.com/tguankheng016/commerce-mono/internal/configurations"
	"github.com/tguankheng016/commerce-mono/internal/data/seeds"
	"github.com/tguankheng016/commerce-mono/internal/identities"
	"github.com/tguankheng016/commerce-mono/internal/users"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/environment"
	"github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"github.com/tguankheng016/commerce-mono/pkg/permissions"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"github.com/tguankheng016/commerce-mono/pkg/security"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				environment.ConfigAppEnv,
				config.InitConfig,
			),
			logging.Module,
			postgres.Module,
			caching.Module,
			security.Module,
			permissions.Module,
			identities.Module,
			users.Module,
			seeds.Module,
			configurations.Module,
			http.Module,
		),
	).Run()
}
