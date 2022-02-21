package db

import (
	"time"

	"github.com/cyverse/QMS/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// AddDefaultUserPlan adds a user plan to the given user. The error plan added to the user is always a basic
// type plan, which doesn't expire.
func AddDefaultUserPlan(db *gorm.DB, username string) (*model.UserPlan, error) {
	wrapMsg := "unable to add the default user plan"
	var err error

	// Get the user ID.
	user, err := GetUser(db, username)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Get the basic plan ID.
	plan, err := GetPlan(db, PlanNameBasic)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Define the user plan.
	effectiveStartDate := time.Now()
	userPlan := model.UserPlan{
		EffectiveStartDate: &effectiveStartDate,
		UserID:             user.ID,
		PlanID:             plan.ID,
	}
	err = db.Select("EffectiveStartDate", "EffectiveEndDate", "UserID", "PlanID").
		Create(&userPlan).Error
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return &userPlan, nil
}

// GetActiveUserPlan retrieves the user plan record that is currently active for the user. The effective start
// date must be before the current date and the effective end date must either be null or after the current date.
// If multiple active user plans exist, the one with the most recent effective start date is used. If no active
// user plans exist for the user then a new one for the basic plan is created.
func GetActiveUserPlan(db *gorm.DB, username string) (*model.UserPlan, error) {
	wrapMsg := "unable to get the active user plan"
	var err error

	// Look up the currently active user plan, adding a new one if it doesn't exist already.
	var userPlan model.UserPlan
	err = db.
		Table("user_plans").
		Joins("JOIN users ON user_plans.user_id=users.id").
		Where("users.user_name=?", username).
		Where(
			db.Where("CURRENT_TIMESTAMP BETWEEN user_plans.effective_start_date AND user_plans.effective_end_date").
				Or("CURRENT_TIMESTAMP > user_plans.effective_start_date AND user_plans.effective_end_date IS NULL"),
		).
		Order("user_plans.effective_start_date desc").
		First(&userPlan).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, wrapMsg)
	} else if err == gorm.ErrRecordNotFound {
		userPlanPtr, err := AddDefaultUserPlan(db, username)
		if err != nil {
			return nil, errors.Wrap(err, wrapMsg)
		}
		userPlan = *userPlanPtr
	}

	return &userPlan, nil
}
