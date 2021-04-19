package configs

import "os"

// env configs
const (
	// EnvServerAddress defines an env variable name
	EnvServerAddress     = "SERVER_ADDR"
	defaultServerAddress = ":3001"
)

// GetServerAddress returns server address
func GetServerAddress() string {
	if addr := os.Getenv(EnvServerAddress); addr != "" {
		return addr
	}
	return defaultServerAddress
}

// InitConfigs initializes yaml configs
func InitConfigs() {
	initServerConfig()
}
