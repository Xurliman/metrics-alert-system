package config

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnvironmentValue(key string) (string, error) {
	if os.Getenv(key) == "" {
		return "", constants.ErrEnvValueMissing
	}
	return os.Getenv(key), nil
}

func GetPort() (string, error) {
	address, err := GetEnvironmentValue("ADDRESS")
	if err != nil {
		return constants.DefaultPort, err
	}

	options := strings.Split(address, ":")
	if len(options) < 2 {
		return constants.DefaultPort, constants.ErrWrongAddress
	}

	port, err := strconv.Atoi(options[1])
	if err != nil {
		return constants.DefaultPort, err
	}

	return ":" + strconv.Itoa(port), nil
}

func GetAppEnv() string {
	appEnv, err := GetEnvironmentValue("APP_ENV")
	if err != nil {
		return "development"
	}
	return appEnv
}

func GetStoreInterval() time.Duration {
	storeIntervalStr, err := GetEnvironmentValue("STORE_INTERVAL")
	if err != nil {
		return time.Second * constants.DefaultStoreInterval
	}

	storeInterval, err := strconv.Atoi(storeIntervalStr)
	if err != nil {
		return time.Second * constants.DefaultStoreInterval
	}

	if storeInterval == 0 {
		return time.Second
	}
	return time.Second * time.Duration(storeInterval)
}

func GetFileStoragePath() string {
	fileStoragePath, err := GetEnvironmentValue("FILE_STORAGE_PATH")
	if err != nil {
		return constants.DefaultFileStoragePath
	}
	return fileStoragePath
}

func GetShouldRestore() bool {
	shouldRestore, err := GetEnvironmentValue("RESTORE")
	if err != nil {
		return constants.DefaultRestore
	}
	return shouldRestore == "true"
}

func GetDBConnection() string {
	dbConnection, err := GetEnvironmentValue("DB_CONNECTION")
	if err != nil {
		return constants.PostgresConnection
	}
	return dbConnection
}

func GetDBHost() string {
	dbHost, err := GetEnvironmentValue("DB_HOST")
	if err != nil {
		return constants.DefaultDBHost
	}
	return dbHost
}

func GetDBPort() int {
	dbPort, err := GetEnvironmentValue("DB_PORT")
	if err != nil {
		return constants.DefaultDBPort
	}
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		utils.Logger.Error("port is not a number switched to default port",
			zap.Int("port", constants.DefaultDBPort),
			zap.Error(err),
		)
		return constants.DefaultDBPort
	}
	return dbPortInt
}

func GetDBName() string {
	dbName, err := GetEnvironmentValue("DB_NAME")
	if err != nil {
		return constants.DefaultDBName
	}
	return dbName
}

func GetDBUsername() string {
	dbUsername, err := GetEnvironmentValue("DB_USERNAME")
	if err != nil {
		return constants.DefaultDBUsername
	}
	return dbUsername
}

func GetDBPassword() string {
	dbPassword, err := GetEnvironmentValue("DB_PASSWORD")
	if err != nil {
		return constants.DefaultDBPassword
	}
	return dbPassword
}

func GetDBSSLMode() string {
	dbSSLMode, err := GetEnvironmentValue("DB_SSL_MODE")
	if err != nil {
		return constants.DefaultDBSSLMode
	}
	return dbSSLMode
}

func GetDBConns() (maxIdleConns int, maxOpenConns int, maxConnLifetime int, maxConnIdleTime int) {
	maxIdleConnsStr, err := GetEnvironmentValue("DB_MAX_IDLE_CONNS")
	if err != nil {
		maxIdleConns = constants.DefaultDBMaxIdleConns
	}
	maxIdleConns, err = strconv.Atoi(maxIdleConnsStr)
	if err != nil {
		maxIdleConns = constants.DefaultDBMaxIdleConns
	}

	maxOpenConnsStr, err := GetEnvironmentValue("DB_MAX_OPEN_CONNS")
	if err != nil {
		maxOpenConns = constants.DefaultDBMaxOpenConns
	}
	maxOpenConns, err = strconv.Atoi(maxOpenConnsStr)
	if err != nil {
		maxOpenConns = constants.DefaultDBMaxOpenConns
	}

	maxConnLifetimeStr, err := GetEnvironmentValue("DB_MAX_CONN_LIFETIME")
	if err != nil {
		maxConnLifetime = constants.DefaultDBMaxConnLifetime
	}
	maxConnLifetime, err = strconv.Atoi(maxConnLifetimeStr)
	if err != nil {
		maxConnLifetime = constants.DefaultDBMaxConnLifetime
	}

	maxConnIdleTimeStr, err := GetEnvironmentValue("DB_MAX_CONN_IDLE_TIME")
	if err != nil {
		maxConnIdleTime = constants.DefaultDBMaxConnIdleTime
	}
	maxConnIdleTime, err = strconv.Atoi(maxConnIdleTimeStr)
	if err != nil {
		maxConnIdleTime = constants.DefaultDBMaxConnIdleTime
	}

	return maxIdleConns, maxOpenConns, maxConnLifetime, maxConnIdleTime
}

func GetDatabaseDSN() string {
	databaseDSN, err := GetEnvironmentValue("DATABASE_DSN")
	if err != nil {
		return ""
	}
	return databaseDSN
}
