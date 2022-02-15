package model

import (
	"gorm.io/gorm"
)

// Usage define the structure for API Usages.
type Usage struct {
	gorm.Model
	ID             *string      `gorm:"json:id" type:"uuid;default:uuid_generate_v1()"`
	Usage          float64      `gorm:"not null"`
	AddedBy        string       `gorm:"json:added_by"`
	LastModifiedBy string       `gorm:"json:last_modified_by"`
	UserPlanID     *string      `gorm:"type:uuid;not null" json:"-"`
	UserPlan       UserPlan     `json:"user_plans"`
	ResourceTypeID *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceType   ResourceType `json:"resource_types"`
}
