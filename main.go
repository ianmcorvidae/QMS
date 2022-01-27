package main

import (
	"github.com/cyverse-de/echo-middleware/v2/log"
	"github.com/cyverse/QMS/server"
	"github.com/sirupsen/logrus"
)

// buildLoggerEntry sets some logging options then returns a logger entry with some custom fields
// for convenience.
func buildLoggerEntry() *logrus.Entry {

	// Enable logging the file name and line number.
	logrus.SetReportCaller(true)

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
	server.Init(logger)
}
