package model

// Users define the structure of an User.
// swagger:model
type Users struct {
	//The id for the User.
	// in: path
	//required: true
	ID *string `gorm:"type:uuid" json:"id"`
	// The Name for the User.
	// in: path
	//required: true
	UserName string `gorm:"not null;unique" json:"username"`
}

func (u *Users) TableName() string {
	return "users"
}
