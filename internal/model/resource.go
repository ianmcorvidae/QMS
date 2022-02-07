package model

// ResourceTypes define the structure for ResourceTypes.
// swagger:model
type ResourceType struct {
	//The id for the Resource.
	// in: path
	//required: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	// The Name for the Resource.
	// in: path
	//required: true
	Name string `gorm:"not null;unique" json:"name"`
	// the Unit/measure for the Resource.
	// in: path
	//required: true
	Unit string `gorm:"not null;unique" json:"unit"`
}
