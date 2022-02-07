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
	//user
	// plans.GET("/", s.GetAllPlans)
	plans.GET("/:plan_id", s.GetPlansForID)
	// should the user be able to view plans using plan ID
	plans.POST("/:plan_name/:description/add", s.AddPlan)
	//Adim
	//add in the body.
	//Post/add
	users := v1.Group("/users")

	//add an user to basic plan.
	users.GET("/:username/plan", s.GetUserPlanDetails)

	// view all the plan details along with usage and quotas
	// users.GET("/:username/quotas", s.GetUserAllQuotas)
	// users.GET("/:username/quotas/:quotaId", s.GetUserQuotaDetails)
	// what if the users have multile palms and wants to look at a particular plan
	// users.GET("/:username/usages", s.GetUserUsageDetails)

	//Admin
	admin := v1.Group("/admin")
	// admin.GET("/users", s.GetAllUsers)
	// admin.GET("/users/:username", s.GetAllUserActivePlans)
	// admin.GET("/quotas", s.GetAllActiveQuotas)

	// admin.PUT("/quotas/:quota_name", s.UpdateQuota)
	admin.GET("/usages", s.GetAllActiveUsage)
	//plans and usages
	admin.POST("/usages", s.UpdateUsages)
	admin.PUT("/user/:user_name/updatePlan/:plan_name", s.UpdateUserQuota)

	//Rersources
	v1.GET("/resources", s.GetAllResources)
	//Route to add admin. when a user is added, the user is automatically assigned Basic plan.
	admin.POST("/:user_name", s.AddUsers)
	//Route to add Resources.
	admin.POST("/resources/:resource_name/:resource_unit/add", s.AddResourceType)
	//TO add Default Quota values for a plan.
	admin.POST("/planQuotaDefault", s.AddPlanQuotaDefault)
	//Update the plan for a particular user.
	admin.PUT("/:user_name/Userplan/:plan_name", s.UpdateUserPlanDetails)

	admin.POST("/users/:user_name/:resource_name/add", s.AddQuota)
	admin.PUT("/user/:user_name/updatePlan/:plan_name", s.UpdateUserQuota)
	admin.POST("/usages/:user_name/:resource_name/add", s.AddUsages)
	//usage/:username/:resource_name/:value

}
