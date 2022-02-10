package controllers

import (
	"fmt"
	"net/http"

	"github.com/cyverse/QMS/internal/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

// swagger:route GET /v1/resource-types resource-types listResourceTypes
//
// List Resource Types
//
// Lists all of the resource types defined in the QMS database.
//
// responses:
//   200: resourceTypeListing
//   400: errorResponse

// swagger:route GET /v1/admin-resource-types admin-resource-types listResourceTypes
//
// List Resource Types
//
// Lists all of the resource types defined in the QMS database.
//
// responses:
//   200: resourceTypeListing
//   400: errorResponse

// ListResourceTypes is the handler for the GET /v1/resource-types and GET /v1/admin/resource-types endpoints.
func (s Server) ListResourceTypes(ctx echo.Context) error {
	data := []model.ResourceType{}
	err := s.GORMDB.Debug().Find(&data).Error
	if err != nil {
		msg := fmt.Sprintf("unable to list resource types: %s", err)
		return model.Error(ctx, msg, http.StatusInternalServerError)
	}
	return model.Success(ctx, data, http.StatusOK)
}

// swagger:route POST /v1/admin/resource-types admin-resource-types addResourceType
//
// Add Resource Type
//
// Adds a new resource type to the QMS database.
//
// responses:
//   200: resourceTypeDetails
//   404: errorResponse

// AddResourceType is the handler for the POST /v1/admin/resource-types endpoint.
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

	// Guard against the case where the ID is specified in the request body, which breaks our duplicate check.
	resourceType.ID = nil

	// Save the resource type.
	result := s.GORMDB.
		Select("ID", "Name", "Unit").
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&resourceType)
	if result.Error != nil {
		msg := fmt.Sprintf("unable to save the resource type: %s", result.Error)
		return model.Error(ctx, msg, http.StatusInternalServerError)
	}

	// If the ID wasn't populated and an error didn't occur then there must have been a conflict.
	if resourceType.ID == nil || *resourceType.ID == "" {
		msg := fmt.Sprintf("a resource type with the given name already exists: %s", resourceType.Name)
		return model.Error(ctx, msg, http.StatusConflict)
	}

	return model.Success(ctx, resourceType, http.StatusOK)
}
