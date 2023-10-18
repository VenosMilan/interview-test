package appconfiguration

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Configuration structure that hold app configuration parameter
type Configuration struct {
	LogDebug       bool
	ServerPort     string
	BinaryFilePath string
}

// NewAppConfiguration constructor for create object configuration
// Config data is from environment variable
func NewAppConfiguration() *Configuration {
	var config = new(Configuration)

	config.ServerPort = os.Getenv("PORT")

	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}

	logDebug, err := strconv.ParseBool(os.Getenv("LOG_DEBUG"))

	if err != nil {
		logDebug = false
	}

	config.LogDebug = logDebug

	if logDebug {
		log.SetLevel(log.DebugLevel)
	}

	config.BinaryFilePath = os.Getenv("BINARY_FILE_PATH")

	if config.BinaryFilePath == "" {
		config.BinaryFilePath = "./records.bin"
	}

	return config
}
