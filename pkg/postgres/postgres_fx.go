package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"postgres_fx",
		postgresProviders,
		postgresInvokes,
	)

	postgresProviders = fx.Options(
		fx.Provide(
			NewPostgresDB,
		),
	)

	postgresInvokes = fx.Options(
		fx.Invoke(RunPostgresDB),
		fx.Invoke(func(db *pgxpool.Pool) error {
			return RunGooseMigration(db)
		}),
	)
)
