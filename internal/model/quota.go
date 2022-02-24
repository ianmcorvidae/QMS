package model

// Quota define the structure for an Api Plan and Quota.
type Quota struct {
	ID             *string      `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	Quota          float64      `json:"quota"`
	AddedBy        string       `gorm:"json:added_by"`
	LastModifiedBy string       `gorm:"json:last_modified_by"`
	UserPlanID     *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceTypeID *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceType   ResourceType `json:"resource_type"`
}
