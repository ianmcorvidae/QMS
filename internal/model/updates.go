package model

import (
	"time"

	"gorm.io/gorm"
)

type UpdateOperations struct {
	ID   string `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v1()"`
	Name string `gorm:"column:name;type:varchar(100)"`
}

func (uo *UpdateOperations) TableName() string {
	return "update_operations"
}

type TrackedMetrics struct {
	Quota string `gorm:"column:quota;type:varchar(100)"`
	Usage string `gorm:"column:usage;type:varchar(100)"`
}

func (tm *TrackedMetrics) TableName() string {
	return "tracked_metrics"
}

type Updates struct {
	gorm.Model
	ID               *string   `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	ValueType        string    `gorm:"not null;unique" json:"value_type"`
	Value            float64   `gorm:"not null;unique" json:"value"`
	UpdatedBy        string    `gorm:"not null" json:"updated_by"`
	EffectiveDate    time.Time `gorm:"json:effective_date;type:date"`
	LastModifiedBy   string    `gorm:"json:last_modified_by;varchar(100)"`
	OpID             string
	UpdateOperations UpdateOperations `gorm:"foreignKey:OpID;references:ID;"`
	ResourceTypeID   string
	ResourceTypes    ResourceTypes `gorm:"foreignKey:ResourceTypeID;references:ID"`
}

func (u *Updates) TableName() string {
	return "updates"
}
