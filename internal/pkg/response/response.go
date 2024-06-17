package response

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const (
	StatusAPISuccess = "SUCCESS"
	StatusAPIError   = "ERROR"
	StatusAPIFailure = "FAILURE"
)

// APIFailureMessage is a default message for failure state.
const APIFailureMessage = "Internal Server Error"

// API represents respond body for HTTP API.
type API struct {
	statusCode int

	Code    int    `json:"code,omitempty"`
	Status  string `json:"status"`
	Entity  string `json:"entity,omitempty"`
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// StatusCode returns status code.
func (a *API) StatusCode() int {
	return a.statusCode
}

// NewHTTPResponse creates a new HTTP Response.
func NewHTTPResponse(entity string) *API {
	return &API{
		Status: StatusAPISuccess,
		Entity: entity,
	}
}

// Success returns response format for success state.
func (a *API) Success(data interface{}, code int, message string) *APISuccess {
	a.statusCode = code
	a.Status = StatusAPISuccess
	a.Message = message
	a.State = formatState(a.Entity, a.Status)
	return &APISuccess{
		API:  a,
		Data: data,
	}
}

// Error returns response format for error state.
func (a *API) Error(code int, message string) *APIError {
	a.statusCode = code
	a.Status = StatusAPIError
	a.Message = message
	a.State = formatState(a.Entity, a.Status)
	return &APIError{
		API: a,
	}
}

// FieldErrors returns response format for validation error.
func (a *API) FieldErrors(err error, code int, message string) *APIError {
	fe := a.Error(code, message)
	fe.Errors = err
	return fe
}

// Failure returns response format for failure state.
func (a *API) Failure(err error, code int) error {
	a.statusCode = code
	a.Status = StatusAPIFailure
	a.Message = APIFailureMessage
	a.State = formatState(a.Entity, a.Status)
	return &APIFailure{
		API:    a,
		causer: err,
	}
}

func (a *API) ErrorWithStatusCode(code int, message string) *APIError {
	a.statusCode = code
	a.Code = code
	a.Status = StatusAPIError
	a.Message = strings.Title(message)
	a.State = formatState(a.Entity, a.Status)
	return &APIError{
		API: a,
	}
}

// APISuccess represents body for API on success.
type APISuccess struct {
	*API
	Meta interface{} `json:"meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// JSON writes response into w.
func (s *APISuccess) JSON(ctx context.Context, w http.ResponseWriter) error {
	return ToJSON(ctx, w, s, s.statusCode)
}

// APIFailure represents body for API on failure. (e.g. Internal Server Error)
type APIFailure struct {
	*API
	causer error
}

// Error implements error interface.
func (f *APIFailure) Error() string {
	b, err := json.Marshal(f) // {"", ""}
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return string(b)
}

// Causer returns error causer.
// The Causer error is needed for logging.
func (f *APIFailure) Causer() error {
	return f.causer
}

// APIError represents response body for API on error.
// e.q: Validation Error, Not Found Error, etc.
type APIError struct {
	*API
	Errors error `json:"errors,omitempty"`
}

// Error implement error interface.
func (e *APIError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return string(b)
}

func formatState(entity string, status string) string {
	status = strings.Title(strings.ToLower(status))
	return entity + status
}

// ToJSON encodes the given data to JSON and write it the given w.
func ToJSON(ctx context.Context, w http.ResponseWriter, data interface{}, status int) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}
	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")
	// Write the status code to the response.
	w.WriteHeader(status)
	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}
	return nil
}

func JsonResponse(w http.ResponseWriter, message string, data interface{}, err interface{}, code int) {
	response := struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
		Errors  interface{} `json:"errors,omitempty"`
	}{
		Code:    code,
		Message: message,
		Data:    data,
		Errors:  err,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}
