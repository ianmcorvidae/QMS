package model

import (
	"gorm.io/gorm"
)

// Usages define the structure for API Usages.
// swagger:model
type Usages struct {
	gorm.Model
	//The id for the plan.
	// in: path
	//required: true
	ID string `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v1()"`
	//The current usage of the resource.
	// in: path
	//required: true
	Usage float64 `gorm:"column:usage;type:numeric"`
	//Added by.
	// in: path
	//required: false
	AddedBy string `gorm:"column:added_by;type:varchar(100)"`
	//The last date the record was modified.
	// in: path
	//required: false
	LastModifiedBy string `gorm:"column:last_modified_by;type:varchar(100)"`
	//The UserPlanID for the Usages.
	// in: path
	//required: true
	UserPlanID string    `gorm:"unique"`
	UserPlans  UserPlans `gorm:"foreignKey:UserPlanID;references:ID;"`
	//The ResourceTypeID for the Usages.
	// in: path
	//required: true
	ResourceTypeID string        `gorm:"unique"`
	ResourceTypes  ResourceTypes `gorm:"foreignKey:ResourceTypeID;references:ID"`
}

func (u *Usages) TableName() string {
	return "usages"
}
