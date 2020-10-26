package configs

import "os"

// This file contains default configs

const (
	defaultServerPort = "3001"

	defaultDatabaseUser     = "root"
	defaultDatabasePassword = "malfplemac"
	defaultDatabaseName     = "otqee"
)

// GetServerPort returns server port
func GetServerPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return port
	}
	return defaultServerPort
}

// GetDatabaseConfig returns database config: (user, password, database_name)
func GetDatabaseConfig() (string, string, string) {
	user := os.Getenv("DATABASE_USER")
	pass := os.Getenv("DATABASE_PASSWORD")
	db := os.Getenv("DATABASE_NAME")
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
