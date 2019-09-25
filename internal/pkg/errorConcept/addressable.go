package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
)

// AddressableErrorConcept encapsulates error concepts which pertain to addressables
type AddressableErrorConcept struct {
	EmptyName addressableEmptyName
	InUse     addressableInUse
	NotFound  addressableNotFound
}

type addressableEmptyName struct{}

func (r addressableEmptyName) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r addressableEmptyName) isA(err error) bool {
	_, ok := err.(errors.ErrEmptyAddressableName)
	return ok
}

type addressableInUse struct{}

func (r addressableInUse) httpErrorCode() int {
	return http.StatusConflict
}

func (r addressableInUse) isA(err error) bool {
	_, ok := err.(errors.ErrAddressableInUse)
	return ok
}

type addressableNotFound struct{}

func (r addressableNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r addressableNotFound) isA(err error) bool {
	_, ok := err.(errors.ErrAddressableNotFound)
	return ok
}
