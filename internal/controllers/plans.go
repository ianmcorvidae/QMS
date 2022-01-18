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
	data := []model.Plans{}
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
	data := model.Plans{}
	err := s.GORMDB.Debug().Where("id=@id", sql.Named("id", plan_id)).Find(&data).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	if data.Name == "" || data.Description == "" {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse("Invalid PlanID", http.StatusInternalServerError))
	}

	return ctx.JSON(http.StatusOK, model.SuccessResponse(data, http.StatusOK))
}
