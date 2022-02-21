package db

import (
	"fmt"

	"github.com/cyverse/QMS/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	PlanNameBasic = "Basic"
)

// GetPlan looks up the plan with the given name.
func GetPlan(db *gorm.DB, planName string) (*model.Plan, error) {
	wrapMsg := fmt.Sprintf("unable to look up plan name '%s'", planName)
	var err error

	plan := model.Plan{Name: planName}
	err = db.First(&plan).Error
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return &plan, nil
}
