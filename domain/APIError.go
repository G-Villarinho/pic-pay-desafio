package domain

import (
	"encoding/json"
	"net/http"
)

var (
	CannotBindPayloadAPIError = NewAPIError(http.StatusUnprocessableEntity, "Invalid Request", "Failed to process the payload")
	InternalServerAPIError    = NewAPIError(http.StatusInternalServerError, "Internal Server Error", "Failed to process the payload")
	SessionNotFoundAPIError   = NewAPIError(http.StatusForbidden, "Authentication Required", "You must be logged in to perform this action. Please log in and try again.")
)

type APIError struct {
	Status int               `json:"status"`
	Title  string            `json:"title"`
	Detail string            `json:"detail"`
	Errors map[string]string `json:"errors,omitempty"`
}

func NewAPIError(status int, title, detail string) *APIError {
	return &APIError{
		Status: status,
		Title:  title,
		Detail: detail,
	}
}

func (e *APIError) WithErrors(errors map[string]string) *APIError {
	e.Errors = errors
	return e
}

func (e *APIError) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *APIError) WriteToResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)

	jsonData, err := e.ToJSON()
	if err != nil {
		return err
	}

	_, err = w.Write(jsonData)
	return err
}
