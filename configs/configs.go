package configs

import "os"

// This file contains default configs

const (
	defaultServerPort = "3001"
)

// GetServerPort returns server port
func GetServerPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return port
	}
	return defaultServerPort
}
