package server

import (
	"fmt"

	"github.com/cyverse-de/echo-middleware/v2/log"
	"github.com/cyverse/QMS/config"
	"github.com/cyverse/QMS/internal/controllers"
	"github.com/cyverse/QMS/internal/db"
)

func Init(logger *log.Logger, spec *config.Specification) {
	e := InitRouter(logger)

	// Establish the database connection.
	logger.Info("establishing the database connection")
	db, gormdb, err := db.Init("postgres", spec.DatabaseURI)
	if err != nil {
		e.Logger.Fatalf("service initialization failed: %s", err.Error())
	}

	s := controllers.Server{
		Router:  e,
		DB:      db,
		GORMDB:  gormdb,
		Service: "qms",
		Title:   "serviceInfo.Title",   //TODO: correct this
		Version: "serviceInfo.Version", //TODO:correct this
	}

	// Register the handlers.
	RegisterHandlers(s)
	e.Logger.Info("starting the service")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", 9000)))
}
