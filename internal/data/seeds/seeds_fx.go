package seeds

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"seeds_fx",
		seedInvoke,
	)

	seedInvoke = fx.Options(
		fx.Invoke(func(ctx context.Context, pool *pgxpool.Pool) error {
			return SeedData(ctx, pool)
		}),
	)
)
