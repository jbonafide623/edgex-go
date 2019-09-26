package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type ErrorConceptType interface {
	httpErrorCode() int
	isA(err error) bool
}

type ExplicitErrorConceptType interface {
	httpErrorCode() int
	isA(err error) bool
	httpError(err error) error
	logMessage(err error) string
}

type handler struct {
	logger logger.LoggingClient
}

type ErrorHandler interface {
	Handle(w http.ResponseWriter, err error, allowableErrors []ErrorConceptType, defaultError ErrorConceptType)
	ExplicitHandle(w http.ResponseWriter, err error, allowableErrors []ExplicitErrorConceptType, defaultError ExplicitErrorConceptType)
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

func (e *handler) ExplicitHandle(w http.ResponseWriter, err error, allowableErrors []ExplicitErrorConceptType, defaultError ExplicitErrorConceptType) {
	var doError = func(err error, message string, errorCode int) {
		e.logger.Error(message)
		http.Error(w, err.Error(), errorCode)
	}

	for key := range allowableErrors {
		if allowableErrors[key].isA(err) {
			ect := allowableErrors[key]
			doError(ect.httpError(err), ect.logMessage(err), ect.httpErrorCode())
			return
		}
	}
	defaultError.logMessage(err)
	doError(defaultError.httpError(err), defaultError.logMessage(err), defaultError.httpErrorCode())
}
