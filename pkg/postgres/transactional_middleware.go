package postgres

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"go.uber.org/zap"
)

const (
	DbContextKey = "Ctx.DbContext.Tx"
)

func SetupTransaction(api huma.API, db *pgxpool.Pool) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		context := ctx.Context()
		tx, err := db.BeginTx(context, pgx.TxOptions{})
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, err.Error())
			return
		}

		ctx = huma.WithValue(ctx, DbContextKey, tx)

		next(ctx)

		if ctx.Status() >= 400 {
			if err := tx.Rollback(context); err != nil {
				logging.Logger.Error("unable to rollback transaction", zap.Error(err))
			}
		} else {
			if err := tx.Commit(context); err != nil && err != pgx.ErrTxClosed {
				huma.WriteErr(api, ctx, http.StatusInternalServerError, "unable to commit transaction", err)
			}
		}
	}
}

func GetTxFromCtx(c context.Context) (pgx.Tx, error) {
	tx, ok := c.Value(DbContextKey).(pgx.Tx)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	return tx, nil
}
