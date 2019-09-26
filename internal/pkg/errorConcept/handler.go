package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type ErrorConceptType interface {
	httpErrorCode() int
	isA(err error) bool
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
