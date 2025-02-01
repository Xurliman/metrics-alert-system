// Package constants defines constant variables and errors used in server side
package constants

const (
	EnvFilePath       = "cmd/server/.env"
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
	LogFilePath       = "storage/logs/server.log"

	StoreIntervalFlag            = "i"
	StoreIntervalFlagDescription = "time interval in seconds to save to the disk"

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

	KeyFlag            = "k"
	KeyFlagDescription = "key to hash the request body"

	DefaultDBMaxIdleConns    = 10
	DefaultDBMaxOpenConns    = 100
	DefaultDBMaxConnLifetime = 0
	DefaultDBMaxConnIdleTime = 8
)
