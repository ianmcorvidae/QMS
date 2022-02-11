package server

import (
	"github.com/cyverse-de/echo-middleware/v2/log"
	"github.com/cyverse-de/echo-middleware/v2/redoc"
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
	e.Use(redoc.Serve(redoc.Opts{Title: "CyVerse Quota Management System"}))

	return e
}

func RegisterHandlers(s controllers.Server) {

	// The base URL acts as a health check endpoint.
	s.Router.GET("/", s.RootHandler)

	// API version 1 endpoints.
	v1 := s.Router.Group("/v1")
	v1.GET("", s.V1RootHandler)

	// Plans.
	plans := v1.Group("/plans")
	plans.GET("", s.GetAllPlans)
	plans.GET("/:plan_id", s.GetPlanByID)
	plans.POST("/:plan_name/:description/add", s.AddPlan)

	// Users.
	users := v1.Group("/users")
	users.GET("/:username/plan", s.GetUserPlanDetails)

	// Resources.
	v1.GET("/resource-types", s.ListResourceTypes)

	// Admin endpoints.
	admin := v1.Group("/admin")

	// Admin usage endpoints.
	admin.GET("/usages", s.GetAllActiveUsage)
	// Updates the usage of a particular user and puts an update record in the update table.
	admin.PUT("/usages", s.UpdateUsages)
	admin.POST("/usages/:user_name/:resource_name", s.AddUsages)
	admin.POST(("/update-op/:update_operation"), s.AddUpdateOperation)

	// Admin user endpoints.
	admin.GET("/users", s.GetAllUsers)
	admin.POST("/:user_name", s.AddUser)
	admin.PUT("/user/:user_name/updatePlan/:plan_name", s.UpdateUserPlan)
	admin.POST("/users/:user_name/:resource_name/:quota_value", s.AddQuota)

	// Admin resource type endpoints.
	adminResourceTypes := admin.Group("/resource-types")
	adminResourceTypes.GET("", s.ListResourceTypes)
	adminResourceTypes.POST("", s.AddResourceType)
	adminResourceTypes.GET("/:resource_type_id", s.GetResourceTypeDetails)
	adminResourceTypes.PUT("/:resource_type_id", s.UpdateResourceType)

	// Admin plan quota default endpoints.
	admin.POST("/planQuotaDefault", s.AddPlanQuotaDefault)
}
