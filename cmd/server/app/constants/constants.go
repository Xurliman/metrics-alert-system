package constants

const (
	EnvFilePath                    = "cmd/server/.env"
	GaugeMetricType                = "gauge"
	CounterMetricType              = "counter"
	TimeKey                        = "timestamp"
	TimestampFormat                = "2006-01-02 15:04:05"
	LogFilePath                    = "storage/logs/server.log"
	DevelopmentMode                = "development"
	ProductionMode                 = "production"
	StoreInterval                  = 300
	FileStoragePath                = "/storage/archive/metrics.json"
	Restore                        = true
	AddressFlag                    = "a"
	StoreIntervalFlag              = "i"
	FileStoragePathFlag            = "f"
	RestoreFlag                    = "r"
	AddressFlagDescription         = "give server host:port (default: localhost:8080)"
	StoreIntervalFlagDescription   = "time interval in seconds to save to the disk"
	FileStoragePathFlagDescription = "path to store archive file"
	RestoreFlagDescription         = "determines whether or not to load previously saved values from the specified file when the server starts"
	DefaultPort                    = ":8080"
)
