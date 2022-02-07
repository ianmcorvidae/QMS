package server

import (
	"github.com/cyverse-de/echo-middleware/v2/log"
	"github.com/cyverse/QMS/internal/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRouter(logger *log.Logger) *echo.Echo {
	// Create the web server.
	e := echo.New()
	// Set a custom logger.
	e.Logger = logger
	// Add middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}

func RegisterHandlers(s controllers.Server) {
	s.Router.GET("/", s.RootHandler)
	// Register the group for API version 1.
	v1 := s.Router.Group("/v1")
	v1.GET("/", s.RootHandler)
	//Plans
	plans := v1.Group("/plans")
	//user
	plans.GET("/:plan_id", s.GetPlansForID)
	// should the user be able to view plans using plan ID
	plans.POST("/:plan_name/:description/add", s.AddPlan)
	//Adim
	//add in the body.
	//Post/add
	users := v1.Group("/users")
	//add an user to basic plan.
	users.GET("/:user_name/plan", s.GetUserPlanDetails)
	//Admin
	admin := v1.Group("/admin")
	admin.GET("/usages", s.GetAllActiveUsage)
	//plans and usages
	admin.POST("/usages", s.UpdateUsages)
	// Admin can Update the UserPlan of the user.
	admin.PUT("/user/:user_name/updatePlan/:plan_name", s.UpdateUserPlan)
	//Route to add User by Admin. When a user is added, the user is automatically assigned Basic plan.
	admin.POST("/:user_name", s.AddUser)
	//Route to add Resources.
	admin.POST("/resources/:resource_name/:resource_unit/add", s.AddResourceType)
	//TO add Default Quota values for a plan.
	admin.POST("/planQuotaDefault", s.AddPlanQuotaDefault)
	//Update the plan for a particular user.
	admin.PUT("/:user_name/:plan_name", s.UpdateUserPlanDetails)
	admin.POST("/users/:user_name/:resource_name/:quota_value", s.AddQuota)
	admin.POST("/usages/:user_name/:resource_name", s.AddUsages)
	//Rersources
	v1.GET("/resources", s.GetAllResources)

}
