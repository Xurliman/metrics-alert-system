package database

import (
	"context"
	"database/sql"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"os"
	"time"
)

func OpenDB(dsn string) error {
	if dsn == "" {
		return constants.ErrDSNEmpty
	}

	db, err := sql.Open(os.Getenv("DB_CONNECTION"), dsn)
	if err != nil {
		return err
	}

	maxIdleConns, maxOpenConns, maxConnLifetime, maxConnIdleTime := config.GetDBConns()
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Duration(maxConnLifetime))
	db.SetConnMaxIdleTime(time.Duration(maxConnIdleTime))
	repositories.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS metrics
		(
			id                  uuid primary key unique not null default uuid_generate_v4(),
			name                text             	    not null,
			metric_type         text             		not null,
			value               double precision            null,
			delta               bigint                      null,
			created_at          timestamp                        default now(),
			updated_at          timestamp,
			deleted_at          timestamp
		);
	`)
	if err != nil {
		return err
	}
	return nil
}
