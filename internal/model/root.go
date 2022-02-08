package model

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RootResponse struct {

	// The name of the service
	Service string `json:"service"`

	// The service title
	Title string `json:"title"`

	// The service version
	Version string `json:"version"`
}

type APIVersionResponse struct {
	RootResponse

	// The API version
	APIVersion string `json:"api_version"`
}

// A wrapper for all response bodies.
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
	Status string      `json:"status"`
}

// Basic success Response
func SuccessResponse(data interface{}, status int) Response {
	return Response{
		Result: data,
		Status: http.StatusText(status),
	}
}

// Success sends a basic success response to the caller.
func Success(ctx echo.Context, data interface{}, status int) error {
	return ctx.JSON(status, SuccessResponse(data, status))
}

// Basic error response
func ErrorResponse(errStr string, status int) Response {
	return Response{
		Error:  errStr,
		Status: http.StatusText(status),
	}
}

// Error sends a basic error response to the caller.
func Error(ctx echo.Context, errStr string, status int) error {
	return ctx.JSON(status, ErrorResponse(errStr, status))
}
