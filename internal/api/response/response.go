// Package response provides standardized API response utilities.
// It implements the response format specified in SPEC.md:
//   - Success: {"status": 200, "result": {...}}
//   - Error: {"status": xxx, "message": "...", "code": "..."}
package response

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse represents a successful API response.
// Status is always 200 for successful responses.
// Result contains the endpoint-specific data.
type SuccessResponse struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
}

// ErrorResponse represents an error API response.
// Status contains the HTTP status code (400, 404, 503, 504, etc.).
// Message contains a human-readable error description in Japanese.
// Code contains a machine-readable error code for client handling.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// JSON writes a JSON response with the given status code.
// It sets the Content-Type header to application/json and encodes the data.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Success writes a successful response with status 200.
// The result parameter is wrapped in SuccessResponse format.
func Success(w http.ResponseWriter, result interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{
		Status: http.StatusOK,
		Result: result,
	})
}

// Error writes an error response with the specified status code.
// message is a Japanese description, code is a machine-readable error code.
func Error(w http.ResponseWriter, status int, message string, code string) {
	JSON(w, status, ErrorResponse{
		Status:  status,
		Message: message,
		Code:    code,
	})
}

// BadRequest writes a 400 Bad Request response.
// Used for invalid parameters or missing required fields.
func BadRequest(w http.ResponseWriter, message string, code string) {
	Error(w, http.StatusBadRequest, message, code)
}

// NotFound writes a 404 Not Found response.
// Used when a requested resource cannot be found.
func NotFound(w http.ResponseWriter, message string, code string) {
	Error(w, http.StatusNotFound, message, code)
}

// InternalServerError writes a 500 Internal Server Error response.
// Used for unexpected server-side errors.
func InternalServerError(w http.ResponseWriter, message string, code string) {
	Error(w, http.StatusInternalServerError, message, code)
}

// ServiceUnavailable writes a 503 Service Unavailable response.
// Used when an external API (Spotify, KKBOX) returns an error.
func ServiceUnavailable(w http.ResponseWriter, message string, code string) {
	Error(w, http.StatusServiceUnavailable, message, code)
}

// GatewayTimeout writes a 504 Gateway Timeout response.
// Used when an external API request times out.
func GatewayTimeout(w http.ResponseWriter, message string, code string) {
	Error(w, http.StatusGatewayTimeout, message, code)
}
