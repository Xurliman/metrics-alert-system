package database

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed database/migrations/*.sql
var embedMigrations embed.FS

func AutoMigrate(db *sql.DB) (err error) {
	goose.SetBaseFS(embedMigrations)
	if err = goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err = goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
