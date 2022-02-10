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
type ErrorResponse struct {

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
type ResponseBodyWrapper struct {

	// The status of the request
	Status string `json:"status"`
}

// General information about the service
// swagger:response rootResponse
type RootResponseWrapper struct {

	// in:body
	Body struct {
		ResponseBodyWrapper

		// The service information
		Result model.RootResponse `json:"result"`
	}
}

// General information about a service API version
//
// swagger:response apiVersionResponse
type APIVersionResponseWrapper struct {

	// in:body
	Body struct {
		ResponseBodyWrapper

		// The API version information
		Result model.APIVersionResponse `json:"result"`
	}
}

// Plan Listing
//
// swagger:response plansResponse
type PlansResponseWrapper struct {

	// in: body
	Body struct {
		ResponseBodyWrapper

		// The list of plans
		Result []model.Plan `json:"result"`
	}
}

// Plan Information
//
// swagger:parameters getPlanByID
type PlanIDParameter struct {

	// The plan identifier
	//
	// in:path
	// required:true
	PlanID string `json:"plan_id"`
}

// Plan Information
//
// swagger:response planResponse
type PlanResponseWrapper struct {

	// in: body
	Body struct {
		ResponseBodyWrapper

		// The plan information
		Result model.Plan `json:"result"`
	}
}

// Users

// User Listing
//
// swagger:response userListing
type UserListingResponseWrapper struct {

	// in: body
	Body struct {
		ResponseBodyWrapper

		// The user listing
		Result []model.User `json:"result"`
	}
}

// Resource Types

// Resource Type Listing
//
// swagger:response resourceTypeListing
type ResourceTypeListingWrapper struct {

	// in: body
	Body struct {
		ResponseBodyWrapper

		// The resource type listing
		Result []model.ResourceType `json:"result"`
	}
}

// Resource Type Details
//
// swagger:response resourceTypeDetails
type ResourceTypeDetailsResponseWrapper struct {

	// in: body
	Body struct {
		ResponseBodyWrapper

		// The resource type information
		Result model.ResourceType `json:"result"`
	}
}

// Parameters for the endpoint used to add resource types.
//
// swagger:parameters addResourceType
type AddResourceTypeParameters struct {

	// The resource type information
	//
	// in: body
	Body model.ResourceType
}
