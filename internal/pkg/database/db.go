package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	appConsts "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"github.com/uptrace/bun/driver/pgdriver"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(options *GormOptions) (*gorm.DB, error) {
	datasource := options.GetDatasource()

	if options.DBName == "" {
		return nil, errors.New("Database name is required in config.json")
	}

	err := createDb(options)

	if err != nil {
		return nil, err
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 3                      // Number of retries (including the initial attempt)

	var gormDb *gorm.DB

	err = backoff.Retry(func() error {

		gormDb, err = gorm.Open(gorm_postgres.Open(datasource), &gorm.Config{})

		if err != nil {
			return errors.Errorf("failed to connect postgres: %v and connection information: %s", err, datasource)
		}

		return nil

	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	return gormDb, err
}

func createDb(options *GormOptions) error {
	datasource := options.GetPostgresDatasource()

	// Create Db If Not Exist
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(datasource)))

	var exists bool
	err := sqldb.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE datname='%s')", options.DBName)).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = sqldb.Exec(fmt.Sprintf("CREATE DATABASE %s", options.DBName))
	if err != nil {
		return err
	}

	defer sqldb.Close()

	return nil
}

func Migrate(gorm *gorm.DB, types ...interface{}) error {
	for _, t := range types {
		err := gorm.AutoMigrate(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func RetrieveTxContext(c echo.Context) (*gorm.DB, error) {
	tx, ok := c.Request().Context().Value(appConsts.TxKey(appConsts.DbContextKey)).(*gorm.DB)
	if !ok {
		return nil, errors.New("Transaction not found in context")
	}

	return tx, nil
}
