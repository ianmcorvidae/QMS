package model

import (
	"time"

	"gorm.io/gorm"
)

// Plans define the structure for an Api Plan.
// swagger:model
type Plans struct {
	//The id for the plan
	// in: path
	//required: true
	ID string `gorm:"column:id;primaryKey,unique;type:uuid"`
	// The Name for the plan
	//required: true
	Name string `gorm:"column:name;type:varchar(100)"`
	// the Description for the plan
	//required: true
	Description string `gorm:"column:description;type:text"`
}

func (p *Plans) TableName() string {
	return "plans"
}

// PlanQuotaDefaults define the structure for an Api Plan and Quota.
// swagger:model
type PlanQuotaDefaults struct {
	gorm.Model
	//The id for the plan.
	// in: path
	//required: true
	ID string `gorm:"column:id;primaryKey;type:uuid"`
	// The QuotaValue for the Quota.
	// in: path
	//required: false
	QuotaValue float64 `gorm:"column:quota_value;type:numeric"`
	// The PlanID for the PlanQuotaDefault.
	//required: false
	PlanID string
	Plans  Plans `gorm:"foreignKey:PlanID;references:ID;"`
	// The ResourceTypeID for the PlanQuotaDefault.
	//required: false
	ResourceTypeID string
	ResourceTypes  ResourceTypes `gorm:"foreignKey:ResourceTypeID;references:ID"`
}

func (pqd *PlanQuotaDefaults) TableName() string {
	return "plan_quota_defaults"
}

type UserPlans struct {
	gorm.Model
	ID                 string    `gorm:"column:id;primaryKey;type:uuid"`
	EffectiveStartDate time.Time `gorm:"column:effective_start_date;type:date"`
	EffectiveEndDate   time.Time `gorm:"column:effective_end_date;type:date"`
	AddedBy            string    `gorm:"column:added_by;type:varchar(100)"`
	LastModifiedBy     string    `gorm:"column:last_modified_by;type:varchar(100)"`
	UserID             string
	Users              Users `gorm:"foreignKey:UserID;references:ID;"`
	PlanID             string
	Plans              Plans `gorm:"foreignKey:PlanID;references:ID;"`
}

func (up *UserPlans) TableName() string {
	return "user_plans"
}
