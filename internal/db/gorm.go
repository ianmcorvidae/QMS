package db

import (
	"database/sql"
	"errors"

	"github.com/cyverse/QMS/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitGORMConnection(db *sql.DB) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	if err != nil {
		return gormDB, errors.New("failed to connect database")
	}
	return gormDB, nil
}

func MigrateTables(db *gorm.DB) error {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.ResourceType{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.Plan{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.PlanQuotaDefault{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.UserPlan{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.Quota{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.UpdateOperation{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.TrackedMetric{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.Update{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.Usage{})
	if err != nil {
		return err
	}
	return nil
}
