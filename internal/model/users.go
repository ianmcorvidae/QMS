package model

// Users define the structure of an User.
// swagger:model
type Users struct {
	//The id for the User.
	// in: path
	//required: true
	ID *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	// The Name for the User.
	// in: path
	//required: true
	UserName string `gorm:"not null;unique" json:"user_name"`
}

func (u *Users) TableName() string {
	return "users"
}
