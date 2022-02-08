package model

import (
	"time"

	"gorm.io/gorm"
)

// Plan
//
// swagger:model
type Plan struct {

	// The plan identifier
	// required: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`

	// The plan name
	// required: true
	Name string `gorm:"not null;unique" json:"name"`

	// A brief description of the plan
	// required: true
	Description string `gorm:"not null" json:"description"`
}

// PlanQuotaDefaults define the structure for an Api Plan and Quota.
type PlanQuotaDefault struct {
	gorm.Model
	ID             *string      `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	PlanID         *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceTypeID *string      `gorm:"type:uuid;not null" json:"-"`
	QuotaValue     float64      `gorm:"not null"`
	Plan           Plan         `json:"plan"`
	ResourceType   ResourceType `json:"resource_type"`
}

// UserPlans define the structure for the API User plans.
type UserPlan struct {
	gorm.Model
	ID                 *string   `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	EffectiveStartDate time.Time `gorm:"json:effective_start_date" type:"date"`
	EffectiveEndDate   time.Time `gorm:"json:effective_end_date" type:"date"`
	AddedBy            string    `gorm:"json:added_by"`
	LastModifiedBy     string    `gorm:"json:last_modified_by"`
	UserID             *string   `gorm:"type:uuid;not null" json:"-"`
	User               User      `json:"users"`
	PlanID             *string   `gorm:"type:uuid;not null" json:"-"`
}
