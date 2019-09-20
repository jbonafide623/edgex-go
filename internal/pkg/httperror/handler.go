package httperror

import (
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

// TODO Figure out the final home for this
type ErrorConceptType interface {
	httpErrorCode() int
	isA(err error) bool
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
