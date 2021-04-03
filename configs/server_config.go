package configs

import (
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
)

// ServerConfig is a struct defining the config from yaml file for this project.
// Please define the config in the config file thoroughly, because there are no null pointer checks
type ServerConfig struct {
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

var serverConfig *ServerConfig

func initServerConfig() {
	logger.GetLogger().Info("init server config")
	serverConfig = &ServerConfig{}

	file, err := os.Open("config.yml")
	if err != nil {
		logger.GetLogger().Fatal("error load config.yml", zap.Error(err))
		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&serverConfig); err != nil {
		logger.GetLogger().Fatal("error decode config", zap.Error(err))
		return
	}
}

// GetServerConfig returns the current config.
// You should never change whatever is returned, because it's a pointer.
func GetServerConfig() *ServerConfig {
	if serverConfig == nil {
		logger.GetLogger().Error("config not initialized")
	}
	return serverConfig
}

// GetDatabaseMySQLDataSourceName returns user, password, database_name combined into a convenient string for mysql
func GetDatabaseMySQLDataSourceName() string {
	c := GetServerConfig()
	return c.Database.Username + ":" +
		c.Database.Password + "@(" +
		c.Database.Address + ")/" +
		c.Database.Name
}
