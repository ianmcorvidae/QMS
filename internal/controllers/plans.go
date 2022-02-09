package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo/v4"
)

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
//   404: errorResponse
func (s Server) GetAllPlans(ctx echo.Context) error {
	data := []model.Plan{}
	err := s.GORMDB.Debug().Find(&data).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, data, http.StatusOK)
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
//   500: errorResponse
func (s Server) GetPlanByID(ctx echo.Context) error {
	plan_id := ctx.Param("plan_id")
	if plan_id == "" {
		return model.Error(ctx, "invalid plan id", http.StatusBadRequest)
	}
	data := model.Plan{}
	err := s.GORMDB.Debug().Where("id=@id", sql.Named("id", plan_id)).Find(&data).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	if data.Name == "" || data.Description == "" {
		msg := fmt.Sprintf("plan id not found: %s", plan_id)
		return model.Error(ctx, msg, http.StatusInternalServerError)
	}

	return model.Success(ctx, data, http.StatusOK)
}

type Plan struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s Server) AddPlan(ctx echo.Context) error {
	planname := ctx.Param("plan_name")
	if planname == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid plan name", http.StatusBadRequest))
	}
	description := ctx.Param("description")
	if description == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid plan description", http.StatusBadRequest))
	}
	var req = model.Plan{Name: planname, Description: description}
	err := s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
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
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("resource Type not found: "+resourceName2, http.StatusInternalServerError))
	}
	storageId := *storage.ID
	var req2 = model.PlanQuotaDefault{PlanID: &planId, ResourceTypeID: &storageId, QuotaValue: storageValue}
	err = s.GORMDB.Debug().Create(&req2).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) UpdateUserPlanDetails(ctx echo.Context) error {
	planname := ctx.Param("plan_name")
	if planname == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid plan name", http.StatusBadRequest))
	}
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid username", http.StatusBadRequest))
	}
	var user = model.User{UserName: username}
	err := s.GORMDB.Debug().Find(&user, "user_name=?", username).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("user not found: "+username, http.StatusInternalServerError))
	}
	userID := *user.ID
	var plan = model.Plan{Name: planname}
	err = s.GORMDB.Debug().Find(&plan, "name=?", planname).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("plan name not found", http.StatusInternalServerError))
	}
	planId := *plan.ID
	var req = model.UserPlan{AddedBy: "Admin", LastModifiedBy: "Admin", UserID: &userID}
	err = s.GORMDB.Debug().Model(&req).Where("user_id=?", userID).Update("plan_id", planId).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) AddQuota(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid username", http.StatusBadRequest))
	}
	resourceName := ctx.Param("resource_name")
	if resourceName == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid resource Name", http.StatusBadRequest))
	}
	quotaValue := ctx.Param("quota_value")
	if quotaValue == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid Quota value", http.StatusBadRequest))
	}
	quotaValueFloat, err := ParseFloat(quotaValue)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("invalid Quota Value", http.StatusInternalServerError))
	}
	var resource = model.ResourceType{Name: resourceName}
	err = s.GORMDB.Debug().Find(&resource, "name=?", resourceName).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("resource Type not found: "+resourceName, http.StatusInternalServerError))
	}
	resourceID := *resource.ID
	var user = model.User{UserName: username}
	err = s.GORMDB.Debug().Find(&user, "user_name=?", username).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("user name Not Found", http.StatusInternalServerError))
	}
	userID := *user.ID
	var userPlan = model.UserPlan{}
	err = s.GORMDB.Debug().Find(&userPlan, "user_id=?", userID).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse("user plan name not found for user: "+username, http.StatusInternalServerError))
	}
	userPlanId := *userPlan.ID
	var req = model.Quota{UserPlanID: &userPlanId, Quota: quotaValueFloat, ResourceTypeID: &resourceID}
	err = s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
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
