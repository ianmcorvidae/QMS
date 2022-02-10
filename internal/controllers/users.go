package controllers

import (
	"net/http"

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
//   404: errorResponse
func (s Server) GetAllUsers(ctx echo.Context) error {
	data := []model.User{}
	err := s.GORMDB.Debug().Find(&data).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(data, http.StatusOK))
}

type PlanDetails struct {
	UserId *string
	Name   string
	Usage  string
	Quota  float64
	Unit   string
}

type QuotaDetails struct {
	PlanName     string
	Quota        float64
	ResourceName string
	Unit         string
}

type UsageDetails struct {
	PlanName     string
	ResourceName string
	Usage        string
	Unit         string
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
	plandata := []PlanDetails{}
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
	users.username=?`, username).Scan(&plandata).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(&plandata, http.StatusOK))
}

func (s Server) AddUser(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid username", http.StatusBadRequest))
	}
	var user = model.User{UserName: username}
	err := s.GORMDB.Debug().Create(&user).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	userID := *user.ID
	var plan = model.Plan{}
	err = s.GORMDB.Debug().Find(&plan, "name=?", "Basic").Error
	if err != nil {
		return model.Error(ctx, "plan name not found.", http.StatusInternalServerError)
	}
	planId := *plan.ID
	var userPlan = model.UserPlan{AddedBy: "Admin", UserID: &userID, PlanID: &planId}
	err = s.GORMDB.Debug().Create(&userPlan).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	storage := "STORAGE"
	cpu := "CPU"
	var storageResource = model.ResourceType{}
	var cpuResource = model.ResourceType{}
	err = s.GORMDB.Debug().Find(&storageResource, "name=?", storage).Error
	if err != nil {
		return model.Error(ctx, "resource not Found: "+storage, http.StatusInternalServerError)
	}
	storageId := *storageResource.ID
	err = s.GORMDB.Debug().Find(&cpuResource, "name=?", cpu).Error
	if err != nil {
		return model.Error(ctx, "resource not found.: "+cpu, http.StatusInternalServerError)
	}
	cpuId := *cpuResource.ID
	userPlanId := *userPlan.ID
	var defaultStorageQuota = model.PlanQuotaDefault{}
	err = s.GORMDB.Debug().Find(&defaultStorageQuota, "resource_type_id=?", storageId).Error
	if err != nil {
		return model.Error(ctx, "default quota not found. for resource type: "+storage,
			http.StatusInternalServerError)
	}
	defaultStorageQuotaValue := defaultStorageQuota.QuotaValue
	var defaultcpuQuota = model.PlanQuotaDefault{}
	err = s.GORMDB.Debug().Find(&defaultcpuQuota, "resource_type_id=?", cpuId).Error
	if err != nil {
		return model.Error(ctx, "default quota not found for resource type: "+cpu,
			http.StatusInternalServerError)
	}
	defaultCPUQuotaValue := defaultcpuQuota.QuotaValue
	var userQuota = []model.Quota{
		{
			AddedBy:        "Admin",
			UserPlanID:     &userPlanId,
			ResourceTypeID: &storageId,
			Quota:          defaultStorageQuotaValue,
		},
		{
			AddedBy:        "Admin",
			UserPlanID:     &userPlanId,
			ResourceTypeID: &cpuId,
			Quota:          defaultCPUQuotaValue,
		},
	}
	err = s.GORMDB.Debug().Create(&userQuota).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) UpdateUserPlan(ctx echo.Context) error {
	planname := ctx.Param("plan_name")
	if planname == "" {
		return model.Error(ctx, "invalid Plan name", http.StatusBadRequest)
	}
	username := ctx.Param("user_name")
	if username == "" {
		return model.Error(ctx, "invalid username", http.StatusBadRequest)
	}
	var user = model.User{UserName: username}
	err := s.GORMDB.Debug().Find(&user, "user_name=?", username).Error
	if err != nil {
		return model.Error(ctx, "user name not found.", http.StatusInternalServerError)
	}
	userID := *user.ID
	var plan = model.Plan{Name: planname}
	err = s.GORMDB.Debug().Find(&plan, "name=?", planname).Error
	if err != nil {
		return model.Error(ctx, "plan name not found.", http.StatusInternalServerError)
	}
	planId := *plan.ID
	var userPlan = model.UserPlan{}
	err = s.GORMDB.Debug().Find(&userPlan, "user_id=?", userID).Error
	if err != nil {
		return model.Error(ctx, "user is not enrolled to a plan", http.StatusInternalServerError)
	}
	userPlanId := *userPlan.ID
	if *userPlan.PlanID == planId {
		return model.Error(ctx, "user cannot be updated with the existing plan: "+planname,
			http.StatusInternalServerError)
	}
	err = s.GORMDB.Debug().
		Model(&userPlan).Where("id=?", userPlanId).
		Update("plan_id", planId).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) AddUsages(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid username", http.StatusBadRequest))
	}
	resourceName := ctx.Param("resource_name")
	if resourceName == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid resource name", http.StatusBadRequest))
	}
	usageValueString := ctx.Param("usage_value")
	if usageValueString == "" {
		return ctx.JSON(http.StatusBadRequest,
			model.ErrorResponse("invalid usage value", http.StatusBadRequest))
	}
	UsageValueFloat, err := ParseFloat(usageValueString)
	if err != nil {
		return model.Error(ctx, "invalid usage value", http.StatusInternalServerError)
	}
	var res = Result{}
	err = s.GORMDB.Table("user_plans").
		Select("user_plans.id, users.user_name,plan_quota_defaults.resource_type_id ").
		Joins("JOIN users on user_plans.user_id=users.id").
		Joins("JOIN plan_quota_defaults on plan_quota_defaults.plan_id=user_plans.plan_id").
		Joins("JOIN resource_types on resource_types.id=plan_quota_defaults.resource_type_id").
		Where("users.user_name=? and resource_types.name=?", username, resourceName).
		Scan(&res).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	var req = model.Usage{
		Usage:          UsageValueFloat,
		AddedBy:        "Admin",
		UserPlanID:     res.ID,
		ResourceTypeID: res.ResourceTypeId,
	}
	err = s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return model.Error(ctx, err.Error(), http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK,
		model.SuccessResponse("Successfully updated Usage for the user: "+username, http.StatusOK))
}
