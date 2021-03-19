package configs

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

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

// Config is a struct defining the config from json file for this project.
// Please define the config in the config file thoroughly, because there are no null pointer checks.
type Config struct {
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Name     string `json:"name"`
	} `json:"database"`
	AllowedOrigins []string `json:"allowed_origins"`
}

var config *Config

// InitConfigs initializes global config struct
func InitConfigs() {
	logger.GetLogger().Info("init configs")
	config = &Config{}

	file, err := os.Open("config.json")
	if err != nil {
		logger.GetLogger().Fatal("error load config.json", zap.Error(err))
		return
	}
	defer file.Close()

	byteVal, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(byteVal, &config); err != nil {
		logger.GetLogger().Fatal("error json unmarshal config", zap.Error(err))
		return
	}
}

// GetConfig returns the current config.
// You should never change whatever is returned, because it's a pointer.
func GetConfig() *Config {
	if config == nil {
		logger.GetLogger().Error("config not initialized")
	}
	return config
}

// GetDatabaseMySQLDataSourceName returns user, password, database_name combined into a convenient string for mysql
func GetDatabaseMySQLDataSourceName() string {
	c := GetConfig()
	return c.Database.Username + ":" +
		c.Database.Password + "@(" +
		c.Database.Address + ")/" +
		c.Database.Name
}
