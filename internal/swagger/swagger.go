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
package swagger

import (
	"github.com/cyverse/QMS/internal/model"
)

// Error Response
//
// swagger:response errorResponse
type errorResponse struct {

	// in: body
	Body struct {

		// A brief description of the error
		Error string `json:"error"`

		// The status of the request
		Status string `json:"status"`
	}
}

// Documentation for the successful response body wrapper. The `Error` field could be included here as well, but it's
// being omitted for now simply because it produces less confusing documentation when the erorr and success response
// bodies are treated separately.
//
// swagger:model
type responseBodyWrapper struct {

	// The status of the request
	Status string `json:"status"`
}

// General information about the service
// swagger:response rootResponse
type rootResponseWrapper struct {

	// in:body
	Body struct {
		responseBodyWrapper

		// The service information
		Result model.RootResponse `json:"result"`
	}
}

// General information about a service API version
// swagger:response apiVersionResponse
type apiVersionResponseWrapper struct {

	// in:body
	Body struct {
		responseBodyWrapper

		// The API version information
		Result model.APIVersionResponse `json:"result"`
	}
}

// Plan Listing
//
// swagger:response plansResponse
type plansResponseWrapper struct {

	// in: body
	Body struct {
		responseBodyWrapper

		// The list of plans
		Result []model.Plan `json:"result"`
	}
}

// Plan Information
//
// swagger:parameters getPlanByID
type plansIDParameter struct {

	// The plan identifier
	//
	// in:path
	// required:true
	PlanID string `json:"plan_id"`
}

// Plan Information
//
// swagger:response planResponse
type planResponseWrapper struct {

	// in: body
	Body struct {
		responseBodyWrapper

		// The plan information
		Result model.Plan `json:"result"`
	}
}

//Users

// User Listing
//
// swagger:response userListing
type userListingResponseWrapper struct {

	// in: body
	Body struct {
		responseBodyWrapper

		// The user listing
		Result []model.User `json:"result"`
	}
}
