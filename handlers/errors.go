package handlers

import (
	"errors"
	"fmt"
)

// APIError define API error when response status is 4xx or 5xx
type APIError struct {
	Code       int         `json:"code"`
	Reason     string      `json:"reason"`
	Message    interface{} `json:"message"`
	TraceId    string      `json:"traceId"`
	ServerTime int64       `json:"serverTime"`
}

// Error return error code and message
func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%d, msg=%s, reason=%s, trace=%s", e.Code, e.Message, e.Reason, e.TraceId)
}

// IsAPIError check if e is an API error
func IsAPIError(e error) bool {
	var APIError *APIError
	ok := errors.As(e, &APIError)
	return ok
}
