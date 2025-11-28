package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response envelope for successful responses
type Response struct {
	Data any   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

// Meta contains pagination information
type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ErrorResponse envelope for error responses
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// Error codes
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeInvalidParam = "INVALID_PARAMETER"
	ErrCodeMissingParam = "MISSING_PARAMETER"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeInternal     = "INTERNAL_ERROR"
)

// respondJSON sends a successful JSON response
func respondJSON(c *gin.Context, status int, data any) {
	c.JSON(status, Response{Data: data})
}

// respondJSONWithMeta sends a successful JSON response with pagination metadata
func respondJSONWithMeta(c *gin.Context, status int, data any, meta *Meta) {
	c.JSON(status, Response{Data: data, Meta: meta})
}

// respondError sends an error JSON response
func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

// respondErrorWithDetails sends an error JSON response with field details
func respondErrorWithDetails(c *gin.Context, status int, code, message string, details map[string]any) {
	c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// respondNotFound sends a 404 not found response
func respondNotFound(c *gin.Context, resource string) {
	respondError(c, http.StatusNotFound, ErrCodeNotFound, resource+" not found")
}

// respondValidationError sends a 400 validation error response
func respondValidationError(c *gin.Context, message string, details map[string]any) {
	respondErrorWithDetails(c, http.StatusBadRequest, ErrCodeValidation, message, details)
}

// respondInvalidParam sends a 400 invalid parameter response
func respondInvalidParam(c *gin.Context, param, message string) {
	respondErrorWithDetails(c, http.StatusBadRequest, ErrCodeInvalidParam, message, map[string]any{
		"parameter": param,
	})
}

// respondMissingParam sends a 400 missing parameter response
func respondMissingParam(c *gin.Context, param string) {
	respondErrorWithDetails(c, http.StatusBadRequest, ErrCodeMissingParam, "Required parameter missing: "+param, map[string]any{
		"parameter": param,
	})
}

// respondInternalError sends a 500 internal error response
func respondInternalError(c *gin.Context) {
	respondError(c, http.StatusInternalServerError, ErrCodeInternal, "An unexpected error occurred")
}

// calculateTotalPages calculates total pages from total count and per page
func calculateTotalPages(total, perPage int) int {
	if perPage <= 0 {
		return 0
	}
	pages := total / perPage
	if total%perPage > 0 {
		pages++
	}
	return pages
}

// Pagination defaults and limits
const (
	DefaultPage    = 1
	DefaultPerPage = 50
	MaxPerPage     = 100
)

// parsePagination extracts and validates pagination params
func parsePagination(c *gin.Context) (page, perPage int, err error) {
	page = DefaultPage
	perPage = DefaultPerPage

	if p := c.Query("page"); p != "" {
		n, err := parseInt(p)
		if err != nil || n < 1 {
			return 0, 0, &paramError{param: "page", message: "must be a positive integer"}
		}
		page = n
	}

	if pp := c.Query("per_page"); pp != "" {
		n, err := parseInt(pp)
		if err != nil || n < 1 {
			return 0, 0, &paramError{param: "per_page", message: "must be a positive integer"}
		}
		perPage = n
		if perPage > MaxPerPage {
			return 0, 0, &paramError{param: "per_page", message: "cannot exceed 100"}
		}
	}

	return page, perPage, nil
}

// paramError represents a parameter validation error
type paramError struct {
	param   string
	message string
}

func (e *paramError) Error() string {
	return e.param + ": " + e.message
}

// parseInt parses a string to int, returns 0 and error on failure
func parseInt(s string) (int, error) {
	if s == "" {
		return 0, &paramError{param: "value", message: "empty string"}
	}
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, &paramError{param: "value", message: "not a number"}
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// calculateOffset calculates SQL offset from page and perPage
func calculateOffset(page, perPage int) int {
	return (page - 1) * perPage
}
