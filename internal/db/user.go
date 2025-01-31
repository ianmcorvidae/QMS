package db

import (
	"github.com/cyverse/QMS/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetUser looks up the user details, adding the user to the database if necessary.
func GetUser(db *gorm.DB, username string) (*model.User, error) {
	wrapMsg := "unable to look up or add the user"
	var err error

	user := model.User{Username: username}
	err = db.Select("ID", "Username").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}},
			UpdateAll: true,
		}).
		Create(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return &user, nil
}
