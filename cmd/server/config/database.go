package config

import (
	"os"
	"strconv"

	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetDBConns() (maxIdleConns int, maxOpenConns int, maxConnLifetime int, maxConnIdleTime int) {
	maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS")
	maxIdleConns, err := strconv.Atoi(maxIdleConnsStr)
	if err != nil {
		maxIdleConns = constants.DefaultDBMaxIdleConns
	}

	maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	maxOpenConns, err = strconv.Atoi(maxOpenConnsStr)
	if err != nil {
		maxOpenConns = constants.DefaultDBMaxOpenConns
	}

	maxConnLifetimeStr := os.Getenv("DB_MAX_CONN_LIFETIME")
	maxConnLifetime, err = strconv.Atoi(maxConnLifetimeStr)
	if err != nil {
		maxConnLifetime = constants.DefaultDBMaxConnLifetime
	}

	maxConnIdleTimeStr := os.Getenv("DB_MAX_CONN_IDLE_TIME")
	maxConnIdleTime, err = strconv.Atoi(maxConnIdleTimeStr)
	if err != nil {
		maxConnIdleTime = constants.DefaultDBMaxConnIdleTime
	}

	return maxIdleConns, maxOpenConns, maxConnLifetime, maxConnIdleTime
}
