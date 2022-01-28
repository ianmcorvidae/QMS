package server

import (
	"github.com/cyverse-de/echo-middleware/v2/log"
	"github.com/cyverse/QMS/internal/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func InitRouter(logger *log.Logger) *echo.Echo {
	// Create the web server.
	e := echo.New()

	// Set a custom logger.
	e.Logger = logger

	// Add middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(redoc.Serve(redoc.Opts{Title: "DE Administrative Requests API Documentation"}))
	// e.Logger.Info("loading service information")
	// serviceInfo, err := getSwaggerServiceInfo()
	// if err != nil {
	// 	e.Logger.Fatal(err, serviceInfo)
	// }
	return e
}

// type SwaggerServiceInfo struct {
// 	Description string `json:"description"`
// 	Title       string `json:"title"`
// 	Version     string `json:"version"`
// }

// // PartialSwagger is a structure that can be used to read just the service information from a Swagger JSON document.
// type PartialSwagger struct {
// 	Info *SwaggerServiceInfo `json:"info"`
// }

// func getSwaggerServiceInfo() (*SwaggerServiceInfo, error) {
// 	wrapMsg := "unable to load the Swagger JSON"

// 	// Open the file containing the Swagger JSON.
// 	file, err := os.Open("swagger.json")
// 	if err != nil {
// 		return nil, errors.Wrap(err, wrapMsg)
// 	}

// 	// Parse the bits we want out of the Swagger JSON.
// 	decoder := json.NewDecoder(file)
// 	partialSwagger := &PartialSwagger{}
// 	err = decoder.Decode(partialSwagger)
// 	if err != nil {
// 		return nil, errors.Wrap(err, wrapMsg)
// 	}
// 	return partialSwagger.Info, nil
// }

func RegisterHandlers(s controllers.Server) {
	s.Router.GET("/", s.RootHandler)

	// Register the group for API version 1.
	v1 := s.Router.Group("/v1")
	v1.GET("/", s.RootHandler)
	//Plans
	plans := v1.Group("/plans")
	// plans.GET("/", s.GetAllPlans)
	plans.GET("/p/:plan_id", s.GetPlansForID)
	// should the user be able to view plans using plan ID
	plans.POST("/", s.AddPlans)
	//Users
	//
	users := v1.Group("/users")
	users.POST("/add", s.AddUsers)

	users.GET("/:username/plan", s.GetUserPlanDetails)

	// view all the plan details along with usage and quotas
	users.GET("/:username/quotas", s.GetUserAllQuotas)
	users.GET("/:username/quotas/:quotaId", s.GetUserQuotaDetails)
	// what if the users have multile palms and wants to look at a particular plan
	// users.GET("/:username/usages", s.GetUserUsageDetails)

	//Admin
	admin := v1.Group("/admin")
	// admin.GET("/users", s.GetAllUsers)
	// admin.GET("/users/:username", s.GetAllUserActivePlans)
	admin.GET("/quotas", s.GetAllActiveQuotas)

	admin.PUT("/quotas/:quotaid", s.UpdateQuota)
	admin.GET("/usages", s.GetAllActiveUsage)
	admin.POST("/usages", s.UpdateUsages)
	admin.PUT("/user/:user_name/updatePlan/:plan_name", s.UpdateUserQuota)

	//Rersources
	v1.GET("/resources", s.GetAllResources)
	v1.POST("/resources/add", s.AddResourceType)
	v1.POST("/planquota/add", s.AddPlanQuotaDefault)
	v1.POST("/Userplan/add", s.AddUserPlanDetails)
	v1.POST("/quotas/add", s.AddQuota)

}
