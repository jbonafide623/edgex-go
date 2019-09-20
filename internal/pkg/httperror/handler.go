package httperror

import (
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

// TODO Figure out the final home for this
type ErrorConceptType interface {
	httpErrorCode() int
	isA(err error) bool
}

type StatusRequestEntityTooLargeErrorConcept struct{}

func (r StatusRequestEntityTooLargeErrorConcept) httpErrorCode() int {
	return http.StatusRequestEntityTooLarge
}

func (r StatusRequestEntityTooLargeErrorConcept) isA(err error) bool {
	switch err.(type) {
	case errors.ErrLimitExceeded:
		return true
	default:
		return false
	}
}

type StatusBadRequestErrorConcept struct{}

func (r StatusBadRequestErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r StatusBadRequestErrorConcept) isA(err error) bool {
	// TODO Is this right?
	return false
}

type StatusInternalServerErrorConcept struct{}

func (r StatusInternalServerErrorConcept) httpErrorCode() int {
	return http.StatusInternalServerError
}

func (r StatusInternalServerErrorConcept) isA(err error) bool {
	// TODO Is this right?
	return false
}

type StatusServiceUnavailableErrorConcept struct{}

func (r StatusServiceUnavailableErrorConcept) httpErrorCode() int {
	return http.StatusServiceUnavailable
}

func (r StatusServiceUnavailableErrorConcept) isA(err error) bool {
	// TODO not sure if this makes the most sense
	return false
}

type StatusNotFoundErrorConcept struct{}

func (r StatusNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r StatusNotFoundErrorConcept) isA(err error) bool {
	return true
}

type StatusConflictErrorConcept struct{}

func (r StatusConflictErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r StatusConflictErrorConcept) isA(err error) bool {
	return true
}

func ToHttpError(w http.ResponseWriter, l logger.LoggingClient, err error, allowableErrors []ErrorConceptType, defaultError ErrorConceptType) {
	// handles error
	var doError = func(errorCode int) {
		message := err.Error()
		l.Error(message)
		http.Error(w, message, errorCode)
	}

	for key := range allowableErrors {
		if allowableErrors[key].isA(err) {
			doError(allowableErrors[key].httpErrorCode())
			return
		}
	}
	doError(defaultError.httpErrorCode())
}
