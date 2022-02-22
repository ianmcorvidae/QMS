package controllers

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"time"

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

type PlanDetails struct {
	UserName string
	UserId   *string
	Name     string
	Usage    string
	Quota    float64
	Unit     string
}

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
	var planData []PlanDetails
	err := s.GORMDB.Debug().Raw(`select plans.name,
	usages.usage,
	quotas.quota,
	resource_types.unit 
	from user_plans 
	join plans on plans.id = user_plans.plan_id 
	join usages on user_plans.id=usages.user_plan_id 
	join quotas on user_plans.id=quotas.user_plan_id 
	join resource_types on resource_types.id=quotas.resource_type_id
	join users on users.id = user_plans.user_id
	where 
	users.username=?`, username).Scan(&planData).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(&planData, http.StatusOK))
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
		AddedBy:            "Admin",
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
			AddedBy:        "Admin",
			UserPlanID:     userPlanId,
			ResourceTypeID: storageId,
			Quota:          defaultStorageQuotaValue,
		},
		{
			AddedBy:        "Admin",
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
		User{UserName: username}
	err := s.GORMDB.Debug().Find(&user, "user_name=?", username).Error
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
}

func (s Server) AddUsages(ctx echo.Context) error {
	var (
		err   error
		usage Usage
	)

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

		// Add the usage record.
		var req = model.Usage{
			Usage:          usage.UsageValue,
			AddedBy:        "Admin",
			UserPlanID:     userPlan.ID,
			ResourceTypeID: resourceType.ID,
		}
		err = tx.Debug().Create(&req).Error
		if err != nil {
			return model.Error(ctx, err.Error(), http.StatusInternalServerError)
		}
		return ctx.JSON(http.StatusOK,
			model.SuccessResponse("Successfully updated Usage for the user: "+usage.Username, http.StatusOK))
	})
	return err
}
