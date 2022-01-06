package model

import (
	"gorm.io/gorm"
)

// Quotas define the structure for an Api Plan and Quota.
// swagger:model
type Quotas struct {
	gorm.Model
	//The id for the Quotas.
	// in: path
	//required: true
	ID string `gorm:"column:id;primaryKey;type:uuid"`
	//The limit on the resource.
	// in: path
	//required: false
	Quota float64 `gorm:"column:quota;type:numeric"`
	//Added by.
	// in: path
	//required: true
	AddedBy string `gorm:"column:added_by;type:varchar(100)"`
	//The last date the record was modified
	// in: path
	//required: false
	LastModifiedBy string `gorm:"column:last_modified_by;type:varchar(100)"`
	//The userPlanId for the Quota.
	// in: path
	//required: true
	UserPlanID string    `gorm:"unique"`
	UserPlans  UserPlans `gorm:"foreignKey:UserPlanID;references:ID;"`
	//The resourceTypeId for the Quota.
	// in: path
	//required: true
	ResourceTypeID string        `gorm:"unique"`
	ResourceTypes  ResourceTypes `gorm:"foreignKey:ResourceTypeID;references:ID"`
}

func (q *Quotas) TableName() string {
	return "quotas"
}
