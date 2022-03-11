package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cyverse-de/echo-middleware/v2/params"
	"github.com/cyverse/QMS/internal/db"
	"github.com/cyverse/QMS/internal/httpmodel"
	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// extractPlanID extracts and validates the plan ID path parameter.
func extractPlanID(ctx echo.Context) (string, error) {
	planID, err := params.ValidatedPathParam(ctx, "plan_id", "uuid_rfc4122")
	if err != nil {
		return "", fmt.Errorf("the plan ID must be a valid UUID")
	}
	return planID, nil
}

// GetAllPlans is the handler for the GET /v1/plans endpoint.
//
// swagger:route GET /v1/plans plans listPlans
//
// List Plans
//
// Lists all of the plans that are currently available.
//
// responses:
//   200: plansResponse
//   500: internalServerErrorResponse
func (s Server) GetAllPlans(ctx echo.Context) error {
	plans, err := db.ListPlans(s.GORMDB)
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, plans, http.StatusOK)
}

// GetPlanByID returns the plan with the given identifier.
//
// swagger:route GET /plans/{plan_id} plans getPlanByID
//
// Get Plan Information
//
// Returns the plan with the given identifier.
//
// responses:
//   200: planResponse
//   400: badRequestResponse
//   500: internalServerErrorResponse
func (s Server) GetPlanByID(ctx echo.Context) error {
	var err error

	// Extract and validate the plan ID.
	planID, err := extractPlanID(ctx)
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusBadRequest)
	}

	// Look up the plan.
	plan, err := db.GetPlanByID(s.GORMDB, planID)
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	if plan == nil {
		msg := fmt.Sprintf("plan ID %s not found", planID)
		return model.Error(ctx, msg, http.StatusNotFound)
	}

	return model.Success(ctx, plan, http.StatusOK)
}

// AddPlan adds a new plan to the database.
//
// swagger:route POST /v1/plans
func (s Server) AddPlan(ctx echo.Context) error {
	var err error

	// Parse and validate the request body.
	var plan httpmodel.NewPlan
	if err = ctx.Bind(&plan); err != nil {
		return model.Error(ctx, err.Error(), http.StatusBadRequest)
	}
	if err = plan.Validate(); err != nil {
		return model.Error(ctx, err.Error(), http.StatusBadRequest)
	}

	// Begin a transaction.
	return s.GORMDB.Transaction(func(tx *gorm.DB) error {

		return nil
	})
}

func (s Server) AddPlanQuotaDefault(ctx echo.Context) error {
	planName := "Basic"
	resourceName1 := "CPU"
	cpuValue := 4.00
	resourceName2 := "STORAGE"
	storageValue := 1000.00
	var plan = model.Plan{Name: planName}
	err := s.GORMDB.Debug().Find(&plan, "name=?", planName).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("plan name not found", http.StatusInternalServerError))
	}
	planId := *plan.ID
	var cpu = model.ResourceType{Name: resourceName1}
	err = s.GORMDB.Debug().Find(&cpu, "name=?", resourceName1).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("resource Type not found: "+resourceName1, http.StatusInternalServerError))
	}
	cpuId := *cpu.ID
	var req = model.PlanQuotaDefault{PlanID: &planId, ResourceTypeID: &cpuId, QuotaValue: cpuValue}
	err = s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	var storage = model.ResourceType{Name: resourceName2}
	err = s.GORMDB.Debug().Find(&storage, "name=?", resourceName2).Error
	if err != nil {
		return model.Error(ctx, "resource Type not found: "+resourceName2, http.StatusInternalServerError)
	}
	storageId := *storage.ID
	var req2 = model.PlanQuotaDefault{PlanID: &planId, ResourceTypeID: &storageId, QuotaValue: storageValue}
	err = s.GORMDB.Debug().Create(&req2).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, "Success", http.StatusOK)
}

func (s Server) AddQuota(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return model.Error(ctx, "invalid username", http.StatusBadRequest)
	}
	resourceName := ctx.Param("resource_name")
	if resourceName == "" {
		return model.Error(ctx, "invalid resource Name", http.StatusBadRequest)
	}
	quotaValue := ctx.Param("quota_value")
	if quotaValue == "" {
		return model.Error(ctx, "invalid Quota value", http.StatusBadRequest)
	}
	quotaValueFloat, err := ParseFloat(quotaValue)
	if err != nil {
		return model.Error(ctx, "invalid Quota Value", http.StatusInternalServerError)
	}
	var resource = model.ResourceType{Name: resourceName}
	err = s.GORMDB.Debug().Find(&resource, "name=?", resourceName).Error
	if err != nil {
		return model.Error(ctx, "resource Type not found: "+resourceName, http.StatusInternalServerError)
	}
	resourceID := *resource.ID
	var user = model.User{Username: username}
	err = s.GORMDB.Debug().Find(&user, "username=?", username).Error
	if err != nil {
		return model.Error(ctx, "user name Not Found", http.StatusInternalServerError)
	}
	userID := *user.ID
	var userPlan = model.UserPlan{}
	err = s.GORMDB.Debug().
		Find(&userPlan, "user_id=?", userID).Error
	if err != nil {
		return model.Error(ctx, "user plan name not found for user: "+username, http.StatusInternalServerError)
	}
	userPlanId := *userPlan.ID
	var quota = model.Quota{
		UserPlanID:     &userPlanId,
		Quota:          quotaValueFloat,
		ResourceTypeID: &resourceID,
	}
	err = s.GORMDB.Debug().
		Create(&quota).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, "Success", http.StatusOK)
}

func ParseFloat(valueString string) (float64, error) {
	valueFloat := 0.0
	if temp, err := strconv.ParseFloat(valueString, 64); err == nil {
		valueFloat = temp
	} else {
		return valueFloat, err
	}
	return valueFloat, nil
}
