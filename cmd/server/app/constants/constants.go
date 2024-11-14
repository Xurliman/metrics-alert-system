package constants

const (
	EnvFilePath       = "cmd/server/.env"
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
	TimeKey           = "timestamp"
	TimestampFormat   = "2006-01-02 15:04:05"
	LogFilePath       = "storage/logs/server.log"
	DevelopmentMode   = "development"
	ProductionMode    = "production"

	StoreIntervalFlag            = "i"
	StoreIntervalFlagDescription = "time interval in seconds to save to the disk"
	DefaultStoreInterval         = 300

	FileStoragePathFlag            = "f"
	FileStoragePathFlagDescription = "path to store archive file"
	DefaultFileStoragePath         = "/storage/archive/metrics.json"

	RestoreFlag            = "r"
	RestoreFlagDescription = "determines whether or not to load previously saved values from the specified file when the server starts"
	DefaultRestore         = true

	AddressFlag            = "a"
	AddressFlagDescription = "give server host:port (default: localhost:8080)"
	DefaultPort            = ":8080"

	DatabaseDSNFlag            = "d"
	DatabaseDSNFlagDescription = "database connection string path"

	DefaultDBHost            = "localhost"
	DefaultDBPort            = 5432
	DefaultDBName            = "postgres"
	DefaultDBUsername        = "postgres"
	DefaultDBPassword        = "postgres"
	DefaultDBSSLMode         = "disable"
	DefaultDBMaxIdleConns    = 10
	DefaultDBMaxOpenConns    = 100
	DefaultDBMaxConnLifetime = 0
	DefaultDBMaxConnIdleTime = 8
	PostgresConnection       = "pgx"
	MysqlConnection          = "mysql"
	MongoConnection          = "mongo"
	SqliteConnection         = "sqlite"
)
