package db

import (
	"fmt"

	"github.com/cyverse/QMS/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GetResourceTypeByName looks up the resource type with the given name.
func GetResourceTypeByName(db *gorm.DB, name string) (*model.ResourceType, error) {
	wrapMsg := fmt.Sprintf("unable to look up resource type '%s'", name)
	var err error

	resourceType := model.ResourceType{Name: name}
	err = db.First(&resourceType).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return &resourceType, nil
}
