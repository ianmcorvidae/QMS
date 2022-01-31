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
	data := []model.Users{}
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

// swagger:route GET /users/{UserName}/plan users listUserPlansByID
// Returns a List all the User Plan Details
// responses:
//   200: UserResponse
//   404: RootResponse
func (s Server) GetUserPlanDetails(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}
	// now := time.Now().Format("2006-01-02")

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

// swagger:route GET /users/{UserName}/quotas users listAllUserQuotaByID
// Returns a Lists the User Quota Details.
// responses:
//   200: UserResponse
//   404: RootResponse
//
func (s Server) GetUserAllQuotas(ctx echo.Context) error {
	resource := ctx.QueryParam("resource")
	resourceFilter := ""
	if resource != "" {
		resourceFilter = `and resource_types.name=?`
	}
	username := ctx.Param("username")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}

	plandata := []QuotaDetails{}
	err := s.GORMDB.Debug().Raw(`select plans.name as plan_name, 
	quotas.quota, 
	resource_types.name as resource_name, 
	resource_types.unit 
	from user_plans	
	join plans on plans.id = user_plans.plan_id		
	join quotas on user_plans.id=quotas.user_plan_id	
	join resource_types on resource_types.id=quotas.resource_type_id	
	join users on users.id = user_plans.user_id	
	where			
	users.username =?`+resourceFilter, username, resource).Scan(&plandata).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(plandata, http.StatusOK))
}

// swagger:route GET /users/{UserName}/quotas/{quotaid} users listUserQuotaByID
// Returns a Lists the User Quota Details.
// responses:
//   200: UserResponse
//   404: RootResponse
//
func (s Server) GetUserQuotaDetails(ctx echo.Context) error {
	quotaid := ctx.Param("quotaId")
	if quotaid == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid quota Id", http.StatusBadRequest))
	}

	username := ctx.Param("username")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}

	quotadata := []QuotaDetails{}
	err := s.GORMDB.Debug().Raw(`select plans.name as plan_name, quotas.quota, resource_types.name as resource_name, resource_types.unit from user_plans	
		join plans on plans.id = user_plans.plan_id		
		join quotas on user_plans.id=quotas.user_plan_id	
		join resource_types on resource_types.id=quotas.resource_type_id	
		join users on users.id = user_plans.user_id	
		where				
		users.username =? and 	
		quotas.id = ?`, username, quotaid).Scan(&quotadata).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(quotadata, http.StatusOK))
}

// swagger:route GET /users/{UserName}/usages users listUserUsageDetailsByID
// Returns a Lists the User Quota Details.
// responses:
//   200: UserResponse
//   404: RootResponse
//
func (s Server) GetUserUsageDetails(ctx echo.Context) error {
	resource := ctx.QueryParam("resource")
	resourceFilter := ""
	if resource != "" {
		resourceFilter = `and resource_types.name=?`
	}
	username := ctx.Param("username")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}

	plandata := []UsageDetails{}
	err := s.GORMDB.Debug().Raw(`select plans.name as plan_name, usages.usage, resource_types.unit, resource_types.name as resource_name from user_plans 	
		join plans on plans.id = user_plans.plan_id 	
		join usages on user_plans.id=usages.user_plan_id 	
		join quotas on user_plans.id=quotas.user_plan_id 	
		join resource_types on resource_types.id=quotas.resource_type_id	
		join users on users.id = user_plans.user_id	
		where 					
		users.username =?`+resourceFilter, username, resource).Scan(&plandata).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(plandata, http.StatusOK))
}

func (s Server) AddUsers(ctx echo.Context) error {
	username := ctx.Param("user_name")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}

	var req = model.Users{UserName: username}
	err := s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}

func (s Server) UpdateUserQuota(ctx echo.Context) error {
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

	var userPlan = model.UserPlans{}
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
	userPlanID := "230d8bd2-7cc5-11ec-90d6-0242ac120003"
	var req = model.Usages{Usage: 5000, AddedBy: "Admin", UserPlanID: userPlanID, ResourceTypeID: "230d8bd2-7cc5-11ec-90d6-0242ac120003"}
	err := s.GORMDB.Debug().Create(&req).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
}
