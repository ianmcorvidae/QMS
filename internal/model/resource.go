package model

// ResourceTypes define the structure for ResourceTypes.
type ResourceType struct {
	ID   *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	Name string  `gorm:"not null;unique" json:"name"`
	Unit string  `gorm:"not null;unique" json:"unit"`
}
