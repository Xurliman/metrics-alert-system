package database

import (
	"database/sql"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"os"
	"time"
)

func OpenDB(dsn string) error {
	db, err := sql.Open(os.Getenv("DB_CONNECTION"), dsn)
	if err != nil {
		return err
	}

	maxIdleConns, maxOpenConns, maxConnLifetime, maxConnIdleTime := config.GetDBConns()
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Duration(maxConnLifetime))
	db.SetConnMaxIdleTime(time.Duration(maxConnIdleTime))

	err = AutoMigrate(db)
	if err != nil {
		return err
	}

	repositories.DB = db
	return nil
}
