package model

import (
	"time"

	"gorm.io/gorm"
)

// Plan defines the structure for an Api Plan.
// swagger:model
type Plan struct {
	//The id for the plan
	// in: path
	//required: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	// The Name for the plan
	//required: true
	Name string `gorm:"not null;unique" json:"name"`
	// The Description for the plan
	//required: true
	Description string `gorm:"not null" json:"description"`
	// PlanQuotaDefaults []PlanQuotaDefaults `json:"quota_defaults"`

}

// PlanQuotaDefaults define the structure for an Api Plan and Quota.
// swagger:model
type PlanQuotaDefaults struct {
	gorm.Model
	ID             *string       `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	PlanID         *string       `gorm:"type:uuid;not null" json:"-"`
	ResourceTypeID *string       `gorm:"type:uuid;not null" json:"-"`
	QuotaValue     float64       `gorm:"not null"`
	ResourceType   ResourceTypes `json:"resource_type"`
}

func (pqd *PlanQuotaDefaults) TableName() string {
	return "plan_quota_defaults"
}

// UserPlans define the structure for the API User plans.
// swagger:model
type UserPlans struct {
	gorm.Model
	//The id for the UserPlans.
	// in: path
	//required: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	//The effective start date of the plan for the user.
	// in: path
	//required: false
	EffectiveStartDate time.Time `gorm:"json:effective_start_date;type:date"`
	//The effective end date of the plan for the user.
	// in: path
	//required: false
	EffectiveEndDate time.Time `gorm:"json:effective_end_date;type:date"`
	//Added by:.
	// in: path
	//required: false
	AddedBy string `gorm:"json:added_by;type:varchar(100)"`
	//The last modified of the plan for the user.
	// in: path
	//required: false
	LastModifiedBy string `gorm:"json:last_modified_by;type:varchar(100)"`
	//The UserID for the user.
	// in: path
	//required: true
	UserID *string `gorm:"type:uuid;not null" json:"-"`
	User   Users   `json:"users"`

	// //The planID for the user.
	// // in: path
	// //required: true
	PlanID *string `gorm:"type:uuid;not null" json:"-"`
}

func (up *UserPlans) TableName() string {
	return "user_plans"
}
