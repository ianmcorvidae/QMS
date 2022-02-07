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
	ID *string `gorm:"json:id;type:uuid;default:uuid_generate_v1()"`
	//The current usage of the resource.
	// in: path
	//required: true
	Usage float64 `gorm:"not null"`
	//Added by.
	// in: path
	//required: false
	AddedBy string `gorm:"json:added_by;type:varchar(100)"`
	//The last date the record was modified.
	// in: path
	//required: false
	LastModifiedBy string `gorm:"json:last_modified_by;type:varchar(100)"`
	//The UserPlanID for the Usages.
	// in: path
	//required: true
	UserPlanID *string   `gorm:"type:uuid;not null" json:"-"`
	UserPlan   UserPlans `json:"user_plans"`
	//The ResourceTypeID for the Usages.
	// in: path
	//required: true
	ResourceTypeID *string       `gorm:"type:uuid;not null" json:"-"`
	ResourceType   ResourceTypes `json:"resource_types"`
}

func (u *Usages) TableName() string {
	return "usages"
}
