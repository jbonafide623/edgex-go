package httperror

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
)

type AddressableEmptyNameErrorConcept struct{}
type AddressableInUseErrorConcept struct{}
type AddressableNotFoundErrorConcept struct{}

func (r AddressableEmptyNameErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r AddressableEmptyNameErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrEmptyAddressableName)
	return ok
}

func (r AddressableInUseErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r AddressableInUseErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrAddressableInUse)
	return ok
}

func (r AddressableNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r AddressableNotFoundErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrAddressableNotFound)
	return ok
}
