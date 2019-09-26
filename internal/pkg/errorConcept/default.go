package errorConcept

import (
	"net/http"
)

// Default represents a fallback error concept only
type DefaultErrorConcept struct {
	BadRequest            badRequest
	RequestEntityTooLarge requestEntityTooLarge
	InternalServerError   internalServerError
	ServiceUnavailable    serviceUnavailable
	NotFound              notFound
	Conflict              conflict
}

// BadRequest represents a default error to be used by the general error handler
type badRequest struct{}

func (r badRequest) httpErrorCode() int {
	return http.StatusBadRequest
}

// isA should not be called considering it is a fallback error concept
func (r badRequest) isA(err error) bool {
	panic("errorConcept.BadRequest is a default error and isA should not be called")
}

func (r badRequest) httpError(err error) error {
	return err
}

func (r badRequest) logMessage(err error) string {
	return err.Error()
}

type requestEntityTooLarge struct{}

func (r requestEntityTooLarge) httpErrorCode() int {
	return http.StatusRequestEntityTooLarge
}

// isA should not be called considering it is a fallback error concept
func (r requestEntityTooLarge) isA(err error) bool {
	panic("errorConcept.RequestEntityTooLarge is a default error and isA should not be called")
}

func (r requestEntityTooLarge) httpError(err error) error {
	return err
}

func (r requestEntityTooLarge) logMessage(err error) string {
	return err.Error()
}

type internalServerError struct{}

func (r internalServerError) httpErrorCode() int {
	return http.StatusInternalServerError
}

// isA should not be called considering it is a fallback error concept
func (r internalServerError) isA(err error) bool {
	panic("errorConcept.InternalServerError is a default error and isA should not be called")
}

func (r internalServerError) httpError(err error) error {
	return err
}

func (r internalServerError) logMessage(err error) string {
	return err.Error()
}

type serviceUnavailable struct{}

func (r serviceUnavailable) httpErrorCode() int {
	return http.StatusServiceUnavailable
}

// isA should not be called considering it is a fallback error concept
func (r serviceUnavailable) isA(err error) bool {
	panic("errorConcept.ServiceUnavailable is a default error and isA should not be called")
}

func (r serviceUnavailable) httpError(err error) error {
	return err
}

func (r serviceUnavailable) logMessage(err error) string {
	return err.Error()
}

type notFound struct{}

func (r notFound) httpErrorCode() int {
	return http.StatusNotFound
}

// isA should not be called considering it is a fallback error concept
func (r notFound) isA(err error) bool {
	panic("errorConcept.NotFound is a default error and isA should not be called")
}

func (r notFound) httpError(err error) error {
	return err
}

func (r notFound) logMessage(err error) string {
	return err.Error()
}

type conflict struct{}

func (r conflict) httpErrorCode() int {
	return http.StatusConflict
}

// isA should not be called considering it is a fallback error concept
func (r conflict) isA(err error) bool {
	panic("errorConcept.Conflict is a default error and isA should not be called")
}

func (r conflict) httpError(err error) error {
	return err
}

func (r conflict) logMessage(err error) string {
	return err.Error()
}
