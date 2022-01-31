package controllers

import (
	"database/sql"
	"net/http"

	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo"
)

// swagger:route GET /plans plans listPlans
// Returns a List all the plans
// responses:
//   200: plansResponse
//   404: RootResponse

func (s Server) GetAllPlans(ctx echo.Context) error {
	data := []model.Plan{}
	err := s.GORMDB.Debug().Find(&data).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(data, http.StatusOK))
}

// swagger:route GET /plans/{PlanID} plans listPlansByID
// Returns a List all the plans
// responses:
//   200: plansResponse
//   500: RootResponse

func (s Server) GetPlansForID(ctx echo.Context) error {
	plan_id := ctx.Param("plan_id")
	if plan_id == "" {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse("Invalid PlanID", http.StatusInternalServerError))
	}
	data := model.Plan{}
	err := s.GORMDB.Debug().Where("id=@id", sql.Named("id", plan_id)).Find(&data).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	if data.Name == "" || data.Description == "" {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse("Invalid PlanID", http.StatusInternalServerError))
	}

	return ctx.JSON(http.StatusOK, model.SuccessResponse(data, http.StatusOK))
}

type Plan struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s Server) AddPlans(ctx echo.Context) error {

	planname := ctx.Param("plan_name")
	if planname == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Plan Name", http.StatusBadRequest))
	}
	description := ctx.Param("description")
	if description == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Plan Name", http.StatusBadRequest))
	}
	var req = model.Plan{Name: planname, Description: description}
	err := s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}
func (s Server) AddResourceType(ctx echo.Context) error {
	// id := "230d8bd2-7cc5-11ec-90d6-0242ac120003"

	name := ctx.Param("resource_name")
	if name == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Resource Name", http.StatusBadRequest))
	}
	unit := ctx.Param("resource_unit")
	if unit == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Resource Name", http.StatusBadRequest))
	}
	var resource_type = model.ResourceTypes{Name: name, Unit: unit}
	err := s.GORMDB.Debug().Create(&resource_type).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}
func (s Server) AddPlanQuotaDefault(ctx echo.Context) error {
	id := "230d8bd2-7cc5-11ec-90d6-0242ac120003"
	planID := "2e146110-7bf1-11ec-90d6-0242ac120003"
	resourceTypeID := "1783e71c-7cb5-11ec-90d6-0242ac120003"
	var req = model.PlanQuotaDefaults{ID: &id, PlanID: &planID, ResourceTypeID: &resourceTypeID, QuotaValue: 1000}
	err := s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) AddUserPlanDetails(ctx echo.Context) error {

	planname := ctx.Param("plan_name")
	if planname == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Plan Name", http.StatusBadRequest))
	}

	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}
	var user = model.Users{UserName: username}
	s.GORMDB.Debug().Find(&user, "user_name=?", username)
	userID := *user.ID

	var plan = model.Plan{Name: planname}
	s.GORMDB.Debug().Find(&plan, "name=?", planname)
	planId := *plan.ID

	var req = model.UserPlans{AddedBy: "Admin", UserID: &userID, PlanID: &planId}
	err := s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) AddQuota(ctx echo.Context) error {
	// id := "6b858690-7cd8-11ec-90d6-0242ac120003"
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}
	resourceName := ctx.Param("resource_name")
	if resourceName == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Plan Name", http.StatusBadRequest))
	}
	var resource = model.ResourceTypes{Name: resourceName}
	s.GORMDB.Debug().Find(&resource, "name=?", resourceName)
	resourceID := *resource.ID

	var user = model.Users{UserName: username}
	s.GORMDB.Debug().Find(&user, "user_name=?", username)
	userID := *user.ID
	var userPlan = model.UserPlans{}
	s.GORMDB.Debug().Find(&userPlan, "user_id=?", userID)
	userPlanId := *userPlan.ID

	var req = model.Quotas{UserPlanID: &userPlanId, Quota: 1000, ResourceTypeID: &resourceID}
	err := s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}
