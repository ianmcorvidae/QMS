package controllers

import (
	"fmt"
	"net/http"

	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo/v4"
)

// AddResourceType is the handler for the POST /v1/admin/resource-types endpoint.
//
// swagger:route POST /v1/admin/resource-types resource-types addResourceType
//
// Add Resource Type
//
// Adds a new resource type to the QMS database.
//
// responses:
//   200: resourceTypeDetails
//   404: errorResponse
func (s Server) AddResourceType(ctx echo.Context) error {
	var err error

	//  Extract and validate the request body.
	var resourceType model.ResourceType
	if err = ctx.Bind(&resourceType); err != nil {
		msg := fmt.Sprintf("invalid request body: %s", err)
		return model.Error(ctx, msg, http.StatusBadRequest)
	}
	if resourceType.Name == "" || resourceType.Unit == "" {
		msg := "the resource type name and unit are both required"
		return model.Error(ctx, msg, http.StatusBadRequest)
	}

	// Save the resource type.
	result := s.GORMDB.Select("ID", "Name", "Unit").Create(&resourceType)
	if result.Error != nil {
		msg := fmt.Sprintf("unable to save the resource type: %s", result.Error)
		return model.Error(ctx, msg, http.StatusInternalServerError)
	}

	return model.Success(ctx, resourceType, http.StatusOK)
}
