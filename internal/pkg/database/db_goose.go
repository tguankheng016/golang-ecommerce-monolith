package database

import (
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func RunGooseMigration(db *gorm.DB) error {
	dir := "../../internal/data/migrations"

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	goose.SetBaseFS(nil)

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	if err := goose.Up(sqlDb, dir); err != nil {
		return err
	}

	return nil
}
