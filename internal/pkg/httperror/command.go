package httperror

import (
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"net/http"
)

type CommandNotFoundErrorConcept struct{}

func (r CommandNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r CommandNotFoundErrorConcept) isA(err error) bool {
	switch err.(type) {
	case errors.ErrItemNotFound:
		return true
	default:
		return false
	}
}
