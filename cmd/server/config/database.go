package config

import (
	"database/sql"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

var DB *sql.DB

func Open(ps string) {
	var connection string

	if ps == "" {
		connection, ps = GetDSN()
	} else {
		connection = "pgx"
	}

	db, err := sql.Open(connection, ps)
	if err != nil {
		log.Fatal(err)
	}

	maxIdleConns, maxOpenConns, maxConnLifetime, maxConnIdleTime := GetDBConns()
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Duration(maxConnLifetime))
	db.SetConnMaxIdleTime(time.Duration(maxConnIdleTime))
	DB = db
}

func GetDSN() (connection, ps string) {
	var (
		host     = GetDBHost()
		port     = GetDBPort()
		username = GetDBUsername()
		password = GetDBPassword()
		dbname   = GetDBName()
		sslmode  = GetDBSSLMode()
	)

	connection = GetDBConnection()
	switch connection {
	case constants.PostgresConnection:
		ps = fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
			host,
			port,
			username,
			password,
			dbname,
			sslmode)
	case constants.MysqlConnection:
		ps = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s",
			username,
			password,
			host,
			port,
			dbname,
		)
	case constants.SqliteConnection:
		ps = fmt.Sprintf("%v.db", dbname)
	}
	return connection, ps
}
