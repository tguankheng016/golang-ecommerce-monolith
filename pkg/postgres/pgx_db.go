package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type myQueryTracer struct {
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	logging.Logger.Info("Executing command", zap.String("sql", data.SQL))

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

// IPgxDbConn is interface to be used by services
// it can be either pool, conn or tx from pgx
type IPgxDbConn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

// NewPostgresDB creates a new PostgreSQL database connection pool using the given options.
// If the database does not exist, it will be created.
// The function will retry up to maxRetries times with an exponential backoff if the connection
// cannot be established.
func NewPostgresDB(ctx context.Context, options *PostgresOptions) (*pgxpool.Pool, error) {
	connStr := options.GetDatasource()

	if options.DBName == "" {
		return nil, fmt.Errorf("database name is required in config.json")
	}

	err := createDb(options)

	if err != nil {
		return nil, err
	}

	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	dbConfig.ConnConfig.Tracer = &myQueryTracer{}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 3

	var dbPool *pgxpool.Pool

	err = backoff.Retry(func() error {
		dbPool, err = pgxpool.NewWithConfig(ctx, dbConfig)
		if err != nil {
			return fmt.Errorf("unable to create connection pool: %w", err)
		}

		return nil

	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1))) // Number of retries (including the initial attempt)

	return dbPool, nil
}

// RunPostgresDB sets up the fx lifecycle hooks for the PostgreSQL database connection pool.
// The OnStart hook pings the database to ensure it is available.
// The OnStop hook closes the database connection pool.
func RunPostgresDB(lc fx.Lifecycle, logger *zap.Logger, db *pgxpool.Pool) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting postgres...")
			if err := db.Ping(ctx); err != nil {
				return err
			}

			logger.Info("postgres connected successfully")

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("disconnecting postgres...")
			db.Close()
			logger.Info("postgres disconnected")

			return nil
		},
	})

	return nil
}

// createDb creates a new PostgreSQL database using the given options.
// If the database does not exist, it will be created.
// The function will return an error if the database cannot be created.
func createDb(options *PostgresOptions) error {
	datasource := options.GetPostgresDatasource()

	// Create Db If Not Exist
	db, err := sql.Open("pgx", datasource)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	defer db.Close()

	var exists bool
	if err = db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE datname='%s')", options.DBName)).Scan(&exists); err != nil {
		return err
	}

	if exists {
		return nil
	}

	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", options.DBName)); err != nil {
		return err
	}

	return nil
}
