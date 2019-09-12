package httperror

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

// HandleGetAllError Handles errors which occur when retrieving all resources for a particular resource by mapping an
// error to an HTTP Response
func HandleGetAllError(w http.ResponseWriter, l logger.LoggingClient, err error) {
	l.Error(err.Error())
	switch err.(type) {
	case errors.ErrLimitExceeded:
		http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleDbErrorWithInternalServerErrorFallback Wraps the DbRetrieveErrorHandler with an InternalServerError status
// code as the default fallback
func HandleDbErrorWithInternalServerErrorFallback(w http.ResponseWriter, l logger.LoggingClient, err error) {
	handleDbRetrieveError(w, l, err, http.StatusInternalServerError)
}

// HandleDbErrorWithBadRequestFallback Wraps the DbRetrieveErrorHandler with a BadRequest status code as the default
// fallback
func HandleDbErrorWithBadRequestFallback(w http.ResponseWriter, l logger.LoggingClient, err error) {
	handleDbRetrieveError(w, l, err, http.StatusBadRequest)
}

// HandleDbErrorWithServiceUnavailableFallback Wraps the DbRetrieveErrorHandler with a ServiceUnavailable status code
// as the default fallback
func HandleDbErrorWithServiceUnavailableFallback(w http.ResponseWriter, l logger.LoggingClient, err error) {
	handleDbRetrieveError(w, l, err, http.StatusServiceUnavailable)
}

// Utility which handles database retrieval errors and uses the given status code as a fallback
func handleDbRetrieveError(w http.ResponseWriter, l logger.LoggingClient, err error, fallbackSc int) {
	l.Error(err.Error())
	if err == db.ErrNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), fallbackSc)
	}
}
