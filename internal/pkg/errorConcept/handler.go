package errorConcept

import (
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type ErrorConceptType interface {
	httpErrorCode() int
	isA(err error) bool
}

type StatusBadRequestErrorConcept struct{}
type StatusRequestEntityTooLargeErrorConcept struct{}
type StatusInternalServerErrorConcept struct{}
type StatusServiceUnavailableErrorConcept struct{}
type StatusNotFoundErrorConcept struct{}
type StatusConflictErrorConcept struct{}
type DatabaseNotFoundErrorConcept struct{}
type DatabaseNotUniqueErrorConcept struct{}
type DuplicateIdentifierErrorConcept struct{}
type ItemNotFoundErrorConcept struct{}

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

func (r StatusBadRequestErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r StatusBadRequestErrorConcept) isA(err error) bool {
	// TODO Is this right?
	return false
}

func (r StatusInternalServerErrorConcept) httpErrorCode() int {
	return http.StatusInternalServerError
}

func (r StatusInternalServerErrorConcept) isA(err error) bool {
	// TODO Is this right?
	return false
}

func (r StatusServiceUnavailableErrorConcept) httpErrorCode() int {
	return http.StatusServiceUnavailable
}

func (r StatusServiceUnavailableErrorConcept) isA(err error) bool {
	// TODO not sure if this makes the most sense
	return false
}

func (r StatusNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r StatusNotFoundErrorConcept) isA(err error) bool {
	return true
}

func (r StatusConflictErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r StatusConflictErrorConcept) isA(err error) bool {
	return true
}

func (r DatabaseNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DatabaseNotFoundErrorConcept) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r DatabaseNotUniqueErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DatabaseNotUniqueErrorConcept) isA(err error) bool {
	return err == db.ErrNotUnique
}

func (r DuplicateIdentifierErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DuplicateIdentifierErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrDuplicateName)
	return ok
}

func (r ItemNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r ItemNotFoundErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrItemNotFound)
	return ok
}

type handler struct {
	logger logger.LoggingClient
}

type ErrorHandler interface {
	Handle(w http.ResponseWriter, err error, allowableErrors []ErrorConceptType, defaultError ErrorConceptType)
}

func NewErrorHandler(l logger.LoggingClient) ErrorHandler {
	h := handler{l}
	return &h
}

func (e *handler) Handle(w http.ResponseWriter, err error, allowableErrors []ErrorConceptType, defaultError ErrorConceptType) {
	// handles error
	var doError = func(errorCode int) {
		message := err.Error()
		e.logger.Error(message)
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
