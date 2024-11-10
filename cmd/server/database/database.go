package database

import (
	"context"
	"database/sql"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"time"
)

func OpenDB(ps string) error {
	if ps == "" {
		return constants.ErrDSNEmpty
	}

	db, err := sql.Open(config.GetDBConnection(), ps)
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
		
		DROP TABLE IF EXISTS metrics;
		DROP TYPE IF EXISTS metric_type CASCADE;
		
		CREATE TYPE metric_type AS ENUM ('gauge', 'counter');
		CREATE TABLE IF NOT EXISTS metrics
		(
			id                  uuid primary key unique not null default uuid_generate_v4(),
			name                text             unique not null,
			metric_type         metric_type             not null,
			value               double precision            null,
			delta               int                         null,
			created_at          timestamp                        default now(),
			updated_at          timestamp,
			deleted_at          timestamp
		);`)
	if err != nil {
		return err
	}
	return nil
}

//
//func GetDSN() (connection, ps string) {
//	var (
//		host     = config.GetDBHost()
//		port     = config.GetDBPort()
//		username = config.GetDBUsername()
//		password = config.GetDBPassword()
//		dbname   = config.GetDBName()
//		sslmode  = config.GetDBSSLMode()
//	)
//
//	connection = config.GetDBConnection()
//	switch connection {
//	case constants.PostgresConnection:
//		ps = fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
//			host,
//			port,
//			username,
//			password,
//			dbname,
//			sslmode)
//	case constants.MysqlConnection:
//		ps = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s",
//			username,
//			password,
//			host,
//			port,
//			dbname,
//		)
//	case constants.SqliteConnection:
//		ps = fmt.Sprintf("%v.db", dbname)
//	}
//	return connection, ps
//}
