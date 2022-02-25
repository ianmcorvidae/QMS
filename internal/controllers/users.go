package controllers

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/cyverse/QMS/internal/db"
	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo/v4"
)

// swagger:route GET /admin/users admin listUsers
//
// List Users
//
// Lists the users registered in the QMS database.
//
// responses:
//   200: userListing
//   500: internalServerErrorResponse
func (s Server) GetAllUsers(ctx echo.Context) error {
	var data []model.User
	err := s.GORMDB.Debug().Find(&data).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(data, http.StatusOK))
}

const (
	ValueTypeQuotas = "quotas"
	ValueTypeUsages = "usages"
)

type Result struct {
	ID             *string
	UserName       string
	ResourceTypeId *string
}

func (s Server) GetUserPlanDetails(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return model.Error(ctx, "invalid username", http.StatusBadRequest)
	}

	// Start a transaction.
	return s.GORMDB.Transaction(func(tx *gorm.DB) error {
		var err error

		// Look up or insert the user.
		user, err := db.GetUser(tx, username)
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}

		// Look up or create the user plan.
		userPlan, err := db.GetActiveUserPlan(tx, user.Username)
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}

		// Retrieve the user plan so that the associations will be loaded.
		result := model.UserPlan{ID: userPlan.ID}
		err = tx.
			Preload("User").
			Preload("Plan").
			Preload("Quotas").
			Preload("Quotas.ResourceType").
			Preload("Usages").
			Preload("Usages.ResourceType").
			First(&result).Error
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}

		// Return the user plan.
		return model.Success(ctx, result, http.StatusOK)
	})
}

func (s Server) AddUser(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return model.Error(ctx, "invalid username", http.StatusBadRequest)
	}

	// Either add the user to the database or look up the existing user information.
	user, err := db.GetUser(s.GORMDB, username)
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}

	userID := user.ID
	var plan = model.Plan{}
	err = s.GORMDB.Debug().Find(&plan, "name=?", "Basic").Error
	if err != nil {
		return model.Error(ctx, "plan name not found.", http.StatusInternalServerError)
	}
	planId := plan.ID
	startDate := time.Now()
	endDate := startDate.AddDate(1, 0, 0)
	var userPlan = model.UserPlan{
		UserID:             userID,
		PlanID:             planId,
		EffectiveStartDate: &startDate,
		EffectiveEndDate:   &endDate,
	}
	err = s.GORMDB.Debug().Create(&userPlan).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	storage := "data.size"
	cpu := "cpu.hours"
	var storageResource = model.ResourceType{}
	var cpuResource = model.ResourceType{}
	err = s.GORMDB.Debug().
		Find(&storageResource, "name=?", storage).Error
	if err != nil {
		return model.Error(ctx, "resource not Found: "+storage, http.StatusInternalServerError)
	}
	storageId := storageResource.ID
	err = s.GORMDB.Debug().
		Find(&cpuResource, "name=?", cpu).Error
	if err != nil {
		return model.Error(ctx, "resource not found.: "+cpu, http.StatusInternalServerError)
	}
	cpuId := cpuResource.ID
	userPlanId := userPlan.ID
	var defaultStorageQuota = model.PlanQuotaDefault{}
	err = s.GORMDB.Debug().
		Find(&defaultStorageQuota, "resource_type_id=?", storageId).Error
	if err != nil {
		return model.Error(ctx, "default quota not found. for resource type: "+storage,
			http.StatusInternalServerError)
	}
	defaultStorageQuotaValue := defaultStorageQuota.QuotaValue
	var defaultCpuQuota = model.PlanQuotaDefault{}
	err = s.GORMDB.Debug().
		Find(&defaultCpuQuota, "resource_type_id=?", cpuId).Error
	if err != nil {
		return model.Error(ctx, "default quota not found for resource type: "+cpu,
			http.StatusInternalServerError)
	}
	defaultCPUQuotaValue := defaultCpuQuota.QuotaValue
	var userQuota = []model.Quota{
		{
			UserPlanID:     userPlanId,
			ResourceTypeID: storageId,
			Quota:          defaultStorageQuotaValue,
		},
		{
			UserPlanID:     userPlanId,
			ResourceTypeID: cpuId,
			Quota:          defaultCPUQuotaValue,
		},
	}
	err = s.GORMDB.Debug().Create(&userQuota).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, "Success", http.StatusOK)
}

func (s Server) UpdateUserPlan(ctx echo.Context) error {
	planName := ctx.Param("plan_name")
	if planName == "" {
		return model.Error(ctx, "invalid Plan name", http.StatusBadRequest)
	}
	username := ctx.Param("user_name")
	if username == "" {
		return model.Error(ctx, "invalid username", http.StatusBadRequest)
	}
	var user = model.
		User{Username: username}
	err := s.GORMDB.Debug().Find(&user, "username=?", username).Error
	if err != nil {
		return model.Error(ctx, "user name not found.", http.StatusInternalServerError)
	}
	userID := *user.ID
	var plan = model.
		Plan{Name: planName}
	err = s.GORMDB.Debug().
		Find(&plan, "name=?", planName).Error
	if err != nil {
		return model.Error(ctx, "plan name not found.", http.StatusInternalServerError)
	}
	planId := *plan.ID
	var userPlan = model.UserPlan{}
	err = s.GORMDB.Debug().
		Find(&userPlan, "user_id=?", userID).Error
	if err != nil {
		return model.Error(ctx, "user is not enrolled to a plan", http.StatusInternalServerError)
	}
	userPlanId := *userPlan.ID
	if *userPlan.PlanID == planId {
		return model.Error(ctx, "user cannot be updated with the existing plan: "+planName,
			http.StatusInternalServerError)
	}
	currentDate := time.Now()
	endDate := currentDate.
		AddDate(1, 0, 0)
	err = s.GORMDB.Debug().
		Model(&userPlan).Where("id=?", userPlanId).
		Update("plan_id", planId).
		Update("effective_start_date", currentDate).
		Update("effective_end_date", endDate).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}

	return model.Success(ctx, "Success", http.StatusOK)
}

type Usage struct {
	Username     string  `json:"username"`
	ResourceName string  `json:"resource_name"`
	UsageValue   float64 `json:"usage_value"`
	UpdateType   string  `json:"update_type"`
}

// AddUsages adds or updates the usage record for a user, plan, and resource type.
func (s Server) AddUsages(ctx echo.Context) error {
	var (
		err   error
		usage Usage
	)

	// Extract and validate the request body.
	if err = ctx.Bind(&usage); err != nil {
		return model.Error(ctx, "invalid request body", http.StatusBadRequest)
	}
	if usage.Username == "" {
		return model.Error(ctx, "invalid username", http.StatusBadRequest)
	}
	if usage.ResourceName == "" {
		return model.Error(ctx, "invalid resource name", http.StatusBadRequest)
	}
	if usage.UsageValue < 0 {
		return model.Error(ctx, "invalid usage value", http.StatusBadRequest)
	}
	if usage.UpdateType == "" {
		return model.Error(ctx, "missing usage update type value", http.StatusBadRequest)
	}

	err = s.GORMDB.Transaction(func(tx *gorm.DB) error {
		// Look up the currently active user plan, adding a default plan if one doesn't exist already.
		userPlan, err := db.GetActiveUserPlan(tx, usage.Username)
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}

		// Look up the resource type.
		resourceType, err := db.GetResourceTypeByName(tx, usage.ResourceName)
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}
		if resourceType == nil {
			return model.Error(ctx, fmt.Sprintf("resource type '%s' does not exist", usage.ResourceName), http.StatusBadRequest)
		}

		// Initialize the new usage record.
		var newUsage = model.Usage{
			Usage:          usage.UsageValue,
			UserPlanID:     userPlan.ID,
			ResourceTypeID: resourceType.ID,
		}

		// Verify that the update operation for the given update type exists.
		updateOperation := model.UpdateOperation{Name: usage.UpdateType}
		err = tx.Debug().First(&updateOperation).Error
		if err == gorm.ErrRecordNotFound {
			return model.Error(ctx, "invalid update type", http.StatusBadRequest)
		}
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}

		// Determine the current usage, which should be zero if the usage record doesn't exist.
		currentUsage := model.Usage{
			UserPlanID:     userPlan.ID,
			ResourceTypeID: resourceType.ID,
		}
		err = tx.Debug().First(&currentUsage).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}

		// Update the new usage based on the values in the request body.
		switch usage.UpdateType {
		case UpdateTypeSet:
			newUsage.Usage = usage.UsageValue
		case UpdateTypeAdd:
			newUsage.Usage = currentUsage.Usage + usage.UsageValue
		default:
			msg := fmt.Sprintf("invalid update type: %s", usage.UpdateType)
			return model.Error(ctx, msg, http.StatusBadRequest)
		}

		// Either add the new usage record or update the existing one.
		err = tx.Debug().Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "user_plan_id",
				},
				{
					Name: "resource_type_id",
				},
			},
			UpdateAll: true,
		}).Create(&newUsage).Error
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}

		// Store an update record in the database.
		update := model.Update{
			Value:             newUsage.Usage,
			ValueType:         ValueTypeUsages,
			ResourceTypeID:    resourceType.ID,
			EffectiveDate:     time.Now(),
			UpdateOperationID: updateOperation.ID,
		}
		err = tx.Debug().Create(&update).Error
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}

		// Return a response to the caller.
		msg := fmt.Sprintf("successfully updated the usage for: %s", usage.Username)
		return model.SuccessMessage(ctx, msg, http.StatusOK)
	})
	return err
}
