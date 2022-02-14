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

// ErrorResponse
//
// Having the same object definition for multiple HTTP response status codes seems to confuse ReDoc, so we're using
// aliases as a workaround.
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

// BadRequestResponse
//
// swagger:response badRequestResponse
type BadRequestResponse struct {
	ErrorResponse
}

// NotFoundResponse
//
// swagger:response notFoundResponse
type NotFoundResponse struct {
	ErrorResponse
}

// ConflictResponse
//
// swagger:response conflictResponse
type ConflictResponse struct {
	ErrorResponse
}

// InternalServerErrorResponse
//
// swagger:response InternalServerErrorResponse
type InternalServerErrorResponse struct {
	ErrorResponse
}

// ResponseBodyWrapper Documentation for the successful response body wrapper. The `Error` field could be included here as well, but it's
// being omitted for now simply because it produces less confusing documentation when the erorr and success response
// bodies are treated separately.
//
// swagger:model
type ResponseBodyWrapper struct {

	// The status of the request
	Status string `json:"status"`
}

// RootResponseWrapper Service Information
//
// swagger:response rootResponse
type RootResponseWrapper struct {

	// in:body
	Body struct {
		ResponseBodyWrapper

		// The service information
		Result model.RootResponse `json:"result"`
	}
}

// APIVersionResponseWrapper
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

// PlansResponseWrapper
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

// PlanIDParameter
//
// swagger:parameters getPlanByID
type PlanIDParameter struct {

	// The plan identifier
	//
	// in:path
	// required:true
	PlanID string `json:"plan_id"`
}

// PlanResponseWrapper
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

// UserListingResponseWrapper
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

// ResourceTypeListingWrapper
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

// ResourceTypeDetailsResponseWrapper
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

// AddResourceTypeParameters Parameters for the endpoint used to add resource types.
//
// swagger:parameters addResourceType
type AddResourceTypeParameters struct {

	// The resource type information
	//
	// in: body
	Body model.ResourceType
}

// GetResourceTypeDetailsParameters Parameters for the endpoint used to get resource type details.
//
// swagger:parameters getResourceTypeDetails
type GetResourceTypeDetailsParameters struct {

	// The resource type ID
	//
	// in: path
	// required: true
	ResourceTypeID string `json:"resource_type_id"`
}

// UpdateResourceTypeParameters Parameters for the endpoint used to update resource types.
//
// swagger:parameters updateResourceType
type UpdateResourceTypeParameters struct {

	// The resource type ID
	//
	// in: path
	// required: true
	ResourceTypeID string `json:"resource_type_id"`

	// The resource type details
	//
	// in: body
	Body model.ResourceType
}
