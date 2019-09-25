package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
)

// Common represents error concepts which apply across core-services
type Common struct {
	RequestEntityTooLarge commonRequestEntityTooLarge
	DuplicateIdentifier   duplicateIdentifier
}

type commonRequestEntityTooLarge struct{}

func (r commonRequestEntityTooLarge) httpErrorCode() int {
	return http.StatusRequestEntityTooLarge
}

// isA should not be called considering it is a fallback error concept
func (r commonRequestEntityTooLarge) isA(err error) bool {
	_, ok := err.(errors.ErrLimitExceeded)
	return ok
}

type duplicateIdentifier struct{}

func (r duplicateIdentifier) httpErrorCode() int {
	return http.StatusConflict
}

// isA should not be called considering it is a fallback error concept
func (r duplicateIdentifier) isA(err error) bool {
	_, ok := err.(errors.ErrDuplicateName)
	return ok
}
