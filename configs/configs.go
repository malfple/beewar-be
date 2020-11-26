package configs

import "os"

// This file contains default configs

const (
	// EnvServerAddress defines an env variable name
	EnvServerAddress     = "SERVER_ADDR"
	defaultServerAddress = ":3001"

	// EnvDatabaseUser defines an env variable name
	EnvDatabaseUser     = "DATABASE_USER"
	defaultDatabaseUser = "root"
	// EnvDatabasePassword defines an env variable name
	EnvDatabasePassword     = "DATABASE_PASSWORD"
	defaultDatabasePassword = "malfplemac"
	// EnvDatabaseName defines an env variable name
	EnvDatabaseName     = "DATABASE_NAME"
	defaultDatabaseName = "otqee"
)

// GetServerAddress returns server address
func GetServerAddress() string {
	if addr := os.Getenv(EnvServerAddress); addr != "" {
		return addr
	}
	return defaultServerAddress
}

// GetDatabaseConfig returns database config: (user, password, database_name)
func GetDatabaseConfig() (string, string, string) {
	user := os.Getenv(EnvDatabaseUser)
	pass := os.Getenv(EnvDatabasePassword)
	db := os.Getenv(EnvDatabaseName)
	if user == "" {
		user = defaultDatabaseUser
	}
	if pass == "" {
		pass = defaultDatabasePassword
	}
	if db == "" {
		db = defaultDatabaseName
	}
	return user, pass, db
}

// GetDatabaseMySQLDataSourceName returns user, password, database_name combined into a convenient string for mysql
func GetDatabaseMySQLDataSourceName() string {
	user, pass, db := GetDatabaseConfig()
	return user + ":" + pass + "@/" + db
}
