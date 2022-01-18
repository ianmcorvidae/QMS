package model

// ResourceTypes define the structure for ResourceTypes.
// swagger:model
type ResourceTypes struct {
	//The id for the Resource.
	// in: path
	//required: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	// The Name for the Resource.
	// in: path
	//required: true
	Name string `gorm:"column:name;type:varchar(100)"`
	// the Unit/measure for the Resource.
	// in: path
	//required: true
	Unit string `gorm:"column:unit;type:varchar(100)"`
}

func (rt *ResourceTypes) TableName() string {
	return "resource_types"
}
