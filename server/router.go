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
	plans.POST("", s.AddPlan)
	plans.GET("/:plan_id", s.GetPlanByID)

	usages := v1.Group("/usages")
	// Admin usage endpoints.
	usages.GET("/:username", s.GetAllUsageOfUser)
	// Add/Updates the usage of a particular user and puts an update record in the update table.
	usages.POST("", s.AddUsages)

	// Users.
	users := v1.Group("/users")
	users.GET("/:username/plan", s.GetUserPlanDetails)
	// Admin user endpoints.
	users.GET("", s.GetAllUsers)
	users.PUT("/:user_name", s.AddUser)
	users.PUT("/:user_name/:plan_name", s.UpdateUserPlan)
	users.POST("/quota", s.AddQuota)
	users.GET("/all_active_users", s.GetAllActiveUserPlans)

	// Admin resource type endpoints.
	adminResourceTypes := v1.Group("/resource-types")
	adminResourceTypes.GET("", s.ListResourceTypes)
	adminResourceTypes.POST("", s.AddResourceType)
	adminResourceTypes.GET("/:resource_type_id", s.GetResourceTypeDetails)
	adminResourceTypes.PUT("/:resource_type_id", s.UpdateResourceType)

	// Admin plan quota default endpoints.
	plans.POST("/quota-defaults", s.AddPlanQuotaDefault)
}
