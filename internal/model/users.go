package model

// Users define the structure of an User.
// swagger:model
type Users struct {
	//The id for the User.
	// in: path
	//required: true
	ID string `gorm:"column:id;primaryKey;type:uuid"`
	// The Name for the User.
	// in: path
	//required: true
	UserName string `gorm:"column:username;type:varchar(100)"`
}

func (u *Users) TableName() string {
	return "users"
}
