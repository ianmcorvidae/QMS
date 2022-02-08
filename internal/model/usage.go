package model

import (
	"gorm.io/gorm"
)

// Usages define the structure for API Usages.
// swagger:model
type Usage struct {
	gorm.Model
	//The id for the plan.
	// in: path
	//required: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	//The current usage of the resource.
	// in: path
	//required: true
	Usage float64 `gorm:"not null"`
	//Added by.
	// in: path
	//required: false
	AddedBy string `gorm:"json:added_by"`
	//The last date the record was modified.
	// in: path
	//required: false
	LastModifiedBy string `gorm:"json:last_modified_by"`
	//The UserPlanID for the Usages.
	// in: path
	//required: true
	UserPlanID *string  `gorm:"type:uuid;not null" json:"-"`
	UserPlan   UserPlan `json:"user_plans"`
	//The ResourceTypeID for the Usages.
	// in: path
	//required: true
	ResourceTypeID *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceType   ResourceType `json:"resource_types"`
}
