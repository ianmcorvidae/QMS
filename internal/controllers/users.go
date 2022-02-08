package controllers

import (
	"net/http"

	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo"
)

// swagger:route GET admin/users adm in listUsers
// Returns a List all the Users By Admin
// responses:ssw
//   200: UserResponse
//   404: RootResponse
func (s Server) GetAllUsers(ctx echo.Context) error {
	data := []model.User{}
	err := s.GORMDB.Debug().Find(&data).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
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

// swagger:route GET /users/{UserName}/plan users listUserPlansByID
// Returns a List all the User Plan Details
// responses:
//   200: UserResponse
//   404: RootResponse
func (s Server) GetUserPlanDetails(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
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
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(&plandata, http.StatusOK))
}

func (s Server) AddUser(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}
	var user = model.User{UserName: username}
	err := s.GORMDB.Debug().Create(&user).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	userID := *user.ID
	var plan = model.Plan{Name: "Basic"}
	s.GORMDB.Debug().Find(&plan, "name=?", "Basic")
	planId := *plan.ID
	var userPlan = model.UserPlan{AddedBy: "Admin", UserID: &userID, PlanID: &planId}
	err = s.GORMDB.Debug().Create(&userPlan).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	storage := "STORAGE"
	cpu := "CPU"
	var storageResource = model.ResourceType{}
	var cpuResource = model.ResourceType{}
	s.GORMDB.Debug().Find(&storageResource, "name=?", storage)
	storageId := *storageResource.ID
	s.GORMDB.Debug().Find(&cpuResource, "name=?", cpu)
	cpuId := *cpuResource.ID
	userPlanId := *userPlan.ID
	var defaultStorageQuota = model.PlanQuotaDefault{}
	s.GORMDB.Debug().Find(&defaultStorageQuota, "resource_type_id=?", storageId)
	defaultStorageQuotaValue := defaultStorageQuota.QuotaValue
	var defaultcpuQuota = model.PlanQuotaDefault{}
	s.GORMDB.Debug().Find(&defaultcpuQuota, "resource_type_id=?", cpuId)
	defaultCPUQuotaValue := defaultcpuQuota.QuotaValue
	var userQuota = []model.Quota{{AddedBy: "Admin", UserPlanID: &userPlanId, ResourceTypeID: &storageId, Quota: defaultStorageQuotaValue}, {AddedBy: "Admin", UserPlanID: &userPlanId, ResourceTypeID: &cpuId, Quota: defaultCPUQuotaValue}}
	err = s.GORMDB.Debug().Create(&userQuota).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) UpdateUserPlan(ctx echo.Context) error {
	planname := ctx.Param("plan_name")
	if planname == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Plan Name", http.StatusBadRequest))
	}
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}
	var user = model.User{UserName: username}
	s.GORMDB.Debug().Find(&user, "user_name=?", username)
	userID := *user.ID
	var plan = model.Plan{Name: planname}
	s.GORMDB.Debug().Find(&plan, "name=?", planname)
	planId := *plan.ID
	var userPlan = model.UserPlan{}
	s.GORMDB.Debug().Find(&userPlan, "user_id=?", userID)
	userPlanId := *userPlan.ID
	if *userPlan.PlanID == planId {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse("User cannot be updated with the existing plan: "+planname, http.StatusInternalServerError))
	}
	errdb := s.GORMDB.Debug().Model(&userPlan).Where("id=?", userPlanId).Update("plan_id", planId).Error
	if errdb != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(errdb.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) AddUsages(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}
	resourceName := ctx.Param("resource_name")
	if resourceName == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Resource Name", http.StatusBadRequest))
	}
	usageValueString := ctx.Param("usage_value")
	if usageValueString == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid Usage Value", http.StatusBadRequest))
	}
	UsageValueFloat := ParseFloat(usageValueString)
	var res = Result{}
	err := s.GORMDB.Table("user_plans").Select("user_plans.id, users.user_name,plan_quota_defaults.resource_type_id ").Joins("JOIN users on user_plans.user_id=users.id").Joins("JOIN plan_quota_defaults on plan_quota_defaults.plan_id=user_plans.plan_id").Joins("JOIN resource_types on resource_types.id=plan_quota_defaults.resource_type_id").Where("users.user_name=? and resource_types.name=?", username, resourceName).Scan(&res).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	var req = model.Usage{Usage: UsageValueFloat, AddedBy: "Admin", UserPlanID: res.ID, ResourceTypeID: res.ResourceTypeId}
	err = s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Successfully updated Usage for the user: "+username, http.StatusOK))
}
