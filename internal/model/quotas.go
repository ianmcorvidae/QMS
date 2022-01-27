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
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	//The limit on the resource.
	// in: path
	//required: false
	Quota float64 `json:"quota"`
	//Added by.
	// in: path
	//required: true
	AddedBy string `gorm:"json:added_by"`
	//The last date the record was modified
	// in: path
	//required: false
	LastModifiedBy string `gorm:"json:last_modified_by"`
	//The userPlanId for the Quota.
	// in: path
	//required: true
	UserPlanID *string   `gorm:"type:uuid;not null" json:"-"`
	UserPlan   UserPlans `json:"user_plans"`
	//The resourceTypeId for the Quota.
	// in: path
	//required: true
	ResourceTypeID *string       `gorm:"type:uuid;not null" json:"-"`
	ResourceType   ResourceTypes `json:"resource_types"`
}

func (q *Quotas) TableName() string {
	return "quotas"
}
