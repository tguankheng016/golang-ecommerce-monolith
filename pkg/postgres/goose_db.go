package postgres

import (
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunGooseMigration(db *pgxpool.Pool) error {
	dir, err := getDataMigrationsPath()
	if err != nil {
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	goose.SetBaseFS(nil)

	sqlDb := stdlib.OpenDBFromPool(db)

	if err := goose.Up(sqlDb, dir); err != nil {
		return err
	}

	return nil
}

func getDataMigrationsPath() (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse up to find the go.mod file
	rootPath := cwd
	for {
		if _, err := os.Stat(filepath.Join(rootPath, "go.mod")); err == nil {
			// Found the go.mod file
			break
		}
		parent := filepath.Dir(rootPath)
		if parent == rootPath {
			return "", err
		}
		rootPath = parent
	}

	// Get the path to the "data/migrations" folder within the project directory
	migrationsPath := filepath.Join(rootPath, "internal/data/migrations")

	return migrationsPath, nil
}
