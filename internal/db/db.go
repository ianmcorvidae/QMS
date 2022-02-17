package db

import (
	"database/sql"

	"github.com/cyverse-de/dbutil"
	"github.com/cyverse/QMS/internal/model"
	_ "github.com/lib/pq"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// addInitialResourceTypes inserts the default resource types in the QMS database if they don't exist already.
func addInitialResourceTypes(gormdb *gorm.DB) error {
	initialResourceTypes := []model.ResourceType{
		{
			Name: "cpu.hours",
			Unit: "cpu hours",
		},
		{
			Name: "data.size",
			Unit: "bytes",
		},
	}

	// Add the resource types.
	for _, rt := range initialResourceTypes {
		err := gormdb.Clauses(clause.OnConflict{DoNothing: true}).Create(&rt).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// addInitialUpdateTypes inserts the default update types in the QMS database if they don't exist already.
func addInitialUpdateOperations(gormdb *gorm.DB) error {
	initialUpdateOperations := []model.UpdateOperation{
		{
			Name: "ADD",
		},
		{
			Name: "SET",
		},
	}

	// Add the update operations.
	for _, op := range initialUpdateOperations {
		err := gormdb.Clauses(clause.OnConflict{DoNothing: true}).Create(&op).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// InitDatabase establishes a database connection and verifies that the database can be reached.
func Init(driverName, databaseURI string) (*sql.DB, *gorm.DB, error) {

	wrapMsg := "unable to initialize the database"
	// Create a database connector to establish the connection for us.
	connector, err := dbutil.NewDefaultConnector("5s")
	if err != nil {
		return nil, nil, errors.Wrap(err, wrapMsg)
	}
	// Establish the database connection.
	conn, err := connector.Connect(driverName, databaseURI)
	if err != nil {
		return nil, nil, errors.Wrap(err, wrapMsg)
	}
	gormdb, err := InitGORMConnection(conn)
	if err != nil {
		return nil, nil, errors.Wrap(err, wrapMsg)
	}
	err = MigrateTables(gormdb)
	if err != nil {
		return nil, nil, errors.Wrap(err, wrapMsg)
	}
	err = addInitialResourceTypes(gormdb)
	if err != nil {
		return nil, nil, errors.Wrap(err, wrapMsg)
	}
	err = addInitialUpdateOperations(gormdb)
	if err != nil {
		return nil, nil, errors.Wrap(err, wrapMsg)
	}

	return conn, gormdb, nil
}
