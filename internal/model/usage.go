package model

import (
	"gorm.io/gorm"
)

// Usage define the structure for API Usages.
type Usage struct {
	gorm.Model
	ID             *string      `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	Usage          float64      `gorm:"not null" json:"usage"`
	AddedBy        string       `gorm:"type:text" json:"added_by"`
	LastModifiedBy string       `gorm:"type:text" json:"last_modified_by"`
	UserPlanID     *string      `gorm:"type:uuid;not null" json:"-"`
	UserPlan       UserPlan     `json:"user_plans"`
	ResourceTypeID *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceType   ResourceType `json:"resource_types"`
}
