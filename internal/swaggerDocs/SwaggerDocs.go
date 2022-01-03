// Package api QMS
//
// Documentation of the QMS Api
//
//     Schemes: http
//     BasePath: /
//     Version: V1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package swaggerDocs

import (
	"github.com/cyverse/QMS/internal/model"
)

// A list of Responses for Plans.
// swagger:response plansResponse
type plansResponseWrapper struct {
	// All the plans
	// in: body
	Body []model.Plans
}

// swagger:parameters listPlansByID
type plansIDParameter struct {
	//in: path
	//required:true
	PlanID string `json:"plan_id`
}

//Users

// A list of Plans for particular UserName.
// Request type listing Users
// swagger:response UserResponse
type UserPlansResponseWrapper struct {
	// All the plans
	// in: body
	Body []model.UserPlans
}

// swagger:parameters listUserPlansByID
type UsersPlansUsernameParameter struct {
	//in: path
	//required:true
	UserName string `json:"username`
}

// swagger:parameters listAllUserQuotaByID
type UserAllQuotaUsernameParameter struct {
	//in: path
	//required:true
	UserName string `json:"username`
}

// swagger:parameters listUserQuotaByID
type UserQuotaUsernameQuotaIdParameter struct {
	//in: path
	//required:true
	UserName string `json:"username`
	//in: path
	//required:true
	QuotaID string `json:"quotaid"`
}

// swagger:parameters listUserUsageDetailsByID
type UserUserUsageDetailsUsernameParameter struct {
	//in: path
	//required:true
	UserName string `json:"username`
}
