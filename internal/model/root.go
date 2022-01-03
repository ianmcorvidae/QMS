package model

import "net/http"

type RootResponse struct {

	// The name of the service.
	Service string `json:"service"`

	// The service title.
	Title string `json:"title"`

	// The service version.
	Version string `json:"version"`
}

// A list of Responses.
// swagger:response RootResponse
type Response struct {
	// The result of the API.
	Result interface{} `json:"result,omitempty"`
	// Error is returned from the API.
	Error string `json:"error,omitempty"`
	// Status of the Request.
	Status string `json:"status"`
}

// Basic success Response
func SuccessResponse(data interface{}, status int) Response {
	return Response{
		Result: data,
		Status: http.StatusText(status),
	}
}

// Basic error response
func ErrorResponse(errStr string, status int) Response {
	return Response{
		Error:  errStr,
		Status: http.StatusText(status),
	}
}
