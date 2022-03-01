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
	err = db.
		Preload("PlanQuotaDefaults").
		Preload("PlanQuotaDefaults.ResourceType").
		First(&plan).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return &plan, nil
}

// GetPlanByID looks up the plan with the given identifier.
func GetPlanByID(db *gorm.DB, planID string) (*model.Plan, error) {
	wrapMsg := fmt.Sprintf("unable to look up plan ID `%s'", planID)
	var err error

	plan := model.Plan{ID: &planID}
	err = db.
		Preload("PlanQuotaDefaults").
		Preload("PlanQuotaDefaults.ResourceType").
		First(&plan).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return &plan, nil
}

// ListPlans lists all of the plans that are currently available.
func ListPlans(db *gorm.DB) ([]*model.Plan, error) {
	wrapMsg := "unable to list plans"
	var err error

	// List the plans.
	var plans []*model.Plan
	err = db.
		Preload("PlanQuotaDefaults").
		Preload("PlanQuotaDefaults.ResourceType").
		Find(&plans).Error
	if err != nil {
		return nil, errors.Wrapf(err, wrapMsg)
	}

	return plans, nil
}

func GetDefaultQuotaForPlan(db *gorm.DB, planID string) ([]model.PlanQuotaDefault, error) {
	wrapMsg := "unable to look up plan name "
	var err error

	var planQuotaDefaults []model.PlanQuotaDefault
	err = db.Model(&planQuotaDefaults).Where("plan_id=?", planID).Scan(&planQuotaDefaults).Error
	//err = db.Find(&planQuotaDefaults).Error
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return planQuotaDefaults, nil
}
