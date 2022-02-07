package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo"
)

type AdminQuotaDetails struct {
	UserName     string
	UserID       string
	PlanName     string
	Quota        float64
	ResourceName string
	Unit         string
}

type AdminUsageDetails struct {
	UserName     string
	UserID       string
	PlanName     string
	ResourceName string
	Usage        float64
	Unit         string
}

type UpdateUsagesReq struct {
	UserName             string  `json:"username"`
	ResourceType         string  `json:"resource_type"`
	UpdateType           string  `json:"update_type"`
	UsageAdjustmentValue float64 `json:"usage_adjustment_value"`
	EffectiveDate        string  `json:"effective_date"`
}

type UpdateQuotaReq struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

// func GetQuota(db *sql.DB, name string)

func (s Server) UpdateQuota(ctx echo.Context) error {
	quotaid := ctx.Param("quotaid")
	if quotaid == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid QuotaID", http.StatusBadRequest))
	}

	req := UpdateQuotaReq{}
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error(), http.StatusBadRequest))
	}

	quota := model.Quotas{}

	err = s.GORMDB.Debug().Where("id = ?", quotaid).Find(&quota).Error
	if err != nil {
		if strings.Contains(err.Error(), "invalid input") {
			return ctx.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error(), http.StatusBadRequest))
		}
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	if *quota.ID != quotaid {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Quota Not Found", http.StatusBadRequest))
	}

	value := req.Value
	if req.Type == "sub" {
		value = -1 * value
	}

	quota.Quota += value

	err = s.GORMDB.Debug().Updates(&quota).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))

}

func (s Server) UpdateUsages(ctx echo.Context) error {
	var (
		err error
		req UpdateUsagesReq
	)

	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error(), http.StatusBadRequest))
	}

	effectivedate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error(), http.StatusBadRequest))
	}

	usageDetails := []model.Usages{}

	if err = s.GORMDB.Debug().Raw(
		`
			SELECT usages.* 
			FROM user_plans
			JOIN usages ON user_plans.id=usages.user_plan_id
			JOIN quotas ON user_plans.id=quotas.user_plan_id
			JOIN resource_types ON resource_types.id=quotas.resource_type_id
			JOIN users ON users.id = user_plans.user_id
			WHERE user_plans.effective_start_date <= CURRENT_DATE
			AND user_plans.effective_end_date >= CURRENT_DATE
			AND users.user_name = ?
			AND resource_types.name = ?
		`,
		req.UserName,
		req.ResourceType,
	).Scan(&usageDetails).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	for _, usagerec := range usageDetails {
		usagerec.UpdatedAt = effectivedate
		value := req.UsageAdjustmentValue

		if req.UpdateType == "sub" {
			value = -1 * value
		}

		usagerec.Usage += value

		err := s.GORMDB.Debug().Updates(&usagerec).Error
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
		}
	}

	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))

}

func (s Server) GetAllActiveQuotas(ctx echo.Context) error {
	var err error

	resource := ctx.QueryParam("resource")
	username := ctx.QueryParam("username")

	plandata := []AdminQuotaDetails{}

	selectedCols := []string{
		"users.user_name",
		"users.id AS user_id",
		"plans.name AS plan_name",
		"quotas.quota",
		"resource_types.name AS resource_name",
		"resource_types.unit",
	}

	tx := s.GORMDB.Debug().Table("user_plans").
		Select(selectedCols).
		Joins("JOIN plans ON user_plans.plan_id = plans.id").
		Joins("JOIN quotas ON user_plans.id = quotas.user_plan_id").
		Joins("JOIN resource_types ON resource_types.id = quotas.resource_type_id").
		Joins("JOIN users ON users.id = user_plans.user_id").
		Where("user_plans.effective_start_date <= CURRENT_DATE").
		Where("user_plans.effective_end_date >= CURRENT_DATE")

	if username != "" {
		tx = tx.Where("users.user_name = ?", username)
	}

	if resource != "" {
		tx = tx.Where("resource_types.name = ?", resource)
	}

	rows, err := tx.Rows()
	defer rows.Close()

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	for rows.Next() {
		var r AdminQuotaDetails
		if err = rows.Scan(
			&r.UserName,
			&r.UserID,
			&r.PlanName,
			&r.Quota,
			&r.ResourceName,
			&r.Unit,
		); err != nil {
			return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
		}
		plandata = append(plandata, r)
	}

	return ctx.JSON(http.StatusOK, model.SuccessResponse(plandata, http.StatusOK))
}

func (s Server) GetAllActiveUsage(ctx echo.Context) error {
	var err error

	resource := ctx.QueryParam("resource")
	resourcefilter := ""
	if resource != "" {
		resourcefilter = ` and resource_types.name = '` + resource + `'`
	}

	username := ctx.QueryParam("username")
	usernamefilter := ""
	if username != "" {
		usernamefilter = ` and users.username = '` + username + `'`
	}

	plandata := []AdminUsageDetails{}

	if err = s.GORMDB.Debug().Raw(
		`
			SELECT users.user_name,
				users.id as user_id,
				plans.name as plan_name,
				usages.usage,
				resource_types.unit,
				resource_types.name as resource_name
			FROM user_plans
			JOIN plans ON plans.id = user_plans.plan_id
			JOIN usages ON user_plans.id=usages.user_plan_id
			JOIN quotas ON user_plans.id=quotas.user_plan_id
			JOIN resource_types ON resource_types.id=quotas.resource_type_id
			JOIN users ON users.id = user_plans.user_id
			WHERE cast(now() as date) between user_plans.effective_start_date and user_plans.effective_end_date
		` + usernamefilter + resourcefilter,
	).Scan(&plandata).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	return ctx.JSON(http.StatusOK, model.SuccessResponse(plandata, http.StatusOK))

}
func (s Server) GetAllUserActivePlans(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse("Invalid UserName", http.StatusBadRequest))
	}
	now := time.Now().Format("2006-01-02")

	plandata := []PlanDetails{}
	err := s.GORMDB.Raw(`select plans.name,usage.usage,quotas.quota,resource_types.unit from
	user_plans
	join plans on plans.id=user_plans.plan_id
	join usages on user_plans.id=usages.user_plan_id
	join quotas on user_plans.id=usages.user_plan_id
	join resource_types on resource_types.id=quotas.resource_type_id
	join users on users.id=user_plans.user_id
	where
	user_plans.effective_start_date<=? and
	user_plans.effective_end_date>=? and
	users.username=?`, now+"T00:00:00", now+"T23:59:59", username).Scan(&plandata).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return ctx.JSON(http.StatusOK, model.SuccessResponse(plandata, http.StatusOK))

}

// func (s Server) AddQuotas(ctx echo.Context) error {
// 	// id := "230d8bd2-7cc5-11ec-90d6-0242ac120003"
// 	planID := "2e146110-7bf1-11ec-90d6-0242ac120003"
// 	resourceTypeID := "1783e71c-7cb5-11ec-90d6-0242ac120003"
// 	var req = model.Quotas{Quota: 1000, AddedBy: "Admin", UserPlanID: "", ResourceTypeID: &resourceTypeID}
// 	err := s.GORMDB.Debug().Create(&req).Error
// 	if err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
// 	}
// 	return ctx.JSON(http.StatusOK, model.SuccessResponse("Success", http.StatusOK))
// }
