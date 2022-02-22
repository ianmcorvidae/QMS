package controllers

import (
	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	UpdateTypeSet = "SET"
	UpdateTypeAdd = "ADD"
)

type UpdateQuotaReq struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

func (s Server) GetAllUsageOfUser(ctx echo.Context) error {
	var err error
	username := ctx.Param("username")
	if username == "" {
		return model.Error(ctx, "invalid username", http.StatusBadRequest)
	}
	var usageData []model.Usage
	usage := s.GORMDB.Debug().
		Joins("JOIN user_plans ON user_plans.id = usages.user_plan_id").
		Joins("JOIN resource_types ON resource_types.id = usages.resource_type_id").
		Joins("JOIN users ON users.id = user_plans.user_id").
		Where("cast(now() as date) between user_plans.effective_start_date and user_plans.effective_end_date")
	if username != "" {
		usage.Where("users.user_name = ?", username)
	}
	if err = usage.Find(&usageData).Error; err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, usageData, http.StatusOK)
}

func (s Server) GetAllActiveUserPlans(ctx echo.Context) error {
	var planData []PlanDetails
	err := s.GORMDB.
		Table("user_plans").
		Joins("join plans ON plans.id=user_plans.plan_id").
		Joins("join usages ON user_plans.id=usages.user_plan_id").
		Joins("join quota ON user_plans.id=usages.user_plan_id").
		Joins("join resource_types ON resource_types.id=usages.resource_type_id").
		Joins("join users ON users.id=user_plans.user_id").
		Where("user_plans.effective_start_date<=cast(now() as date)").
		Where("user_plans.effective_end_date>=cast(now() as date)").
		Where("usages.resource_type_id=quota.resource_type_id").
		Scan(&planData).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, planData, http.StatusOK)
}

func (s Server) AddUpdateOperation(ctx echo.Context) error {
	updateOperationName := ctx.Param("update_operation")
	if updateOperationName == "" {
		return model.Error(ctx, "invalid update operation", http.StatusBadRequest)
	}
	var updateOperation = model.UpdateOperation{Name: updateOperationName}
	err := s.GORMDB.Debug().Create(&updateOperation).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, "Success", http.StatusOK)
}
