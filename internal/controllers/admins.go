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
	var user model.User
	err = s.GORMDB.Where("username=?", username).Find(&user).Error
	if err != nil {
		return model.Error(ctx, "user not found", http.StatusInternalServerError)
	}
	var userPlan model.UserPlan
	err = s.GORMDB.
		Preload("User").
		Preload("Plan").
		Preload("Usages").
		Preload("Usages.ResourceType").
		Where("user_id=?", user.ID).
		Find(&userPlan).Error

	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, userPlan.Usages, http.StatusOK)
}

func (s Server) GetAllActiveUserPlans(ctx echo.Context) error {
	var userPlan []model.UserPlan
	err := s.GORMDB.
		Preload("User").
		Preload("Plan").
		Preload("Plan.PlanQuotaDefaults").
		Preload("Plan.PlanQuotaDefaults.ResourceType").
		Preload("Quotas").
		Preload("Quotas.ResourceType").
		Preload("Usages").
		Preload("Usages.ResourceType").
		Where(
			s.GORMDB.Where("CURRENT_TIMESTAMP BETWEEN user_plans.effective_start_date AND user_plans.effective_end_date").
				Or("CURRENT_TIMESTAMP > user_plans.effective_start_date AND user_plans.effective_end_date IS NULL")).
		Find(&userPlan).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return model.Success(ctx, userPlan, http.StatusOK)
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
