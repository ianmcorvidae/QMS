package main

import (
	"github.com/cyverse-de/echo-middleware/v2/log"
	"github.com/cyverse/QMS/config"
	"github.com/cyverse/QMS/server"
	"github.com/sirupsen/logrus"
)

// buildLoggerEntry sets some logging options then returns a logger entry with some custom fields
// for convenience.
func buildLoggerEntry() *logrus.Entry {

	// Set the logging format to JSON because that's what Echo's middleware uses.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Return the custom log entry.
	return logrus.WithFields(logrus.Fields{
		"service": "qms",
		"art-id":  "qms",
		"group":   "org.cyverse",
	})
}

func main() {
	logger := log.NewLogger(buildLoggerEntry())

	// Load the configuration.
	spec, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("unable to load the configuration: %s", err.Error())
	}

	// Initialize the server.
	server.Init(logger, spec)
}
