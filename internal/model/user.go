package model

// Users define the structure of an User.
// swagger:model
type User struct {
	//The id for the User.
	// in: path
	//required: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	// The Name for the User.
	// in: path
	//required: true
	UserName string `gorm:"not null;unique" json:"user_name"`
}
