package common

import (
	"os"
	"strconv"
)

const (
	envVarIsProd            = "GOPHR_PROD"
	envVarHTTPPort          = "PORT"
	envVarIsStaging         = "GOPHR_STAGING"
	fallbackHTTPPort        = 3000
	fallbackDatabaseAddress = "gophr.dev"
	envVarDatabaseIPAddress = "GOPHR_DB_PORT_9042_TCP_ADDR"
)

func ReadEnvDatabaseAddress() string {
	databaseIPAddress := os.Getenv(envVarDatabaseIPAddress)
	if len(databaseIPAddress) > 0 {
		return databaseIPAddress
	}

	return fallbackDatabaseAddress
}

func ReadEnvIsDev() bool {
	if len(os.Getenv(envVarIsProd)) > 0 || len(os.Getenv(envVarIsStaging)) > 0 {
		return false
	}

	return true
}

func ReadEnvHTTPPort() int {
	var (
		port    int
		portStr = os.Getenv(envVarHTTPPort)
	)

	if len(portStr) == 0 {
		port = fallbackHTTPPort
	} else if portNum, err := strconv.Atoi(portStr); err == nil {
		port = portNum
	} else {
		port = fallbackHTTPPort
	}

	return port
}
