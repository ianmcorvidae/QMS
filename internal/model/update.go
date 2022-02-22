package model

import (
	"time"

	"gorm.io/gorm"
)

// UpdateOperation defines the structure of an available update operation in the QMS database.
//
// swagger:model
type UpdateOperation struct {
	// The update operation ID
	//
	// required: true
	// readOnly: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	// The update operation name
	//
	// required: true
	Name string `gorm:"type:text;not null;unique" json:"name"`
}

type TrackedMetric struct {
	Quota string `gorm:"not null" json:"quota"`
	Usage string `gorm:"not null" json:"usage"`
}

type Update struct {
	gorm.Model
	ID                *string      `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	ValueType         string       `json:"value_type"`
	Value             float64      `gorm:"not null" json:"value"`
	UpdatedBy         string       `gorm:"not null" json:"updated_by"`
	EffectiveDate     time.Time    `gorm:"json:effective_date" type:"date"`
	LastModifiedBy    string       `gorm:"json:last_modified_by"`
	UpdateOperationID *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceTypeID    *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceType      ResourceType `json:"resource_types"`
}
