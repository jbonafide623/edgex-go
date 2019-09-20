package httperror

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
)

type DuplicateIdentifierErrorConcept struct{}

func (r DuplicateIdentifierErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DuplicateIdentifierErrorConcept) isA(err error) bool {
	switch err.(type) {
	case errors.ErrDuplicateName:
		return true
	default:
		return false
	}
}

type EmptyAddressableNameErrorConcept struct{}

func (r EmptyAddressableNameErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r EmptyAddressableNameErrorConcept) isA(err error) bool {
	switch err.(type) {
	case errors.ErrEmptyAddressableName:
		return true
	default:
		return false
	}
}

type AddressableNotFoundErrorConcept struct{}

func (r AddressableNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r AddressableNotFoundErrorConcept) isA(err error) bool {
	switch err.(type) {
	case errors.ErrAddressableNotFound:
		return true
	default:
		return false
	}
}

type DatabaseNotFoundErrorConcept struct{}

func (r DatabaseNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DatabaseNotFoundErrorConcept) isA(err error) bool {
	return err == db.ErrNotFound
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

type StatusInternalServerErrorConcept struct{}

func (r StatusInternalServerErrorConcept) httpErrorCode() int {
	return http.StatusInternalServerError
}

func (r StatusInternalServerErrorConcept) isA(err error) bool {
	return false
}

type AddressableInUseErrorConcept struct{}

func (r AddressableInUseErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r AddressableInUseErrorConcept) isA(err error) bool {
	switch err.(type) {
	case errors.ErrAddressableInUse:
		return true
	default:
		return false
	}
}

type AddressableNotFoundByNameErrorConcept struct{}

func (r AddressableNotFoundByNameErrorConcept) httpErrorCode() int {
	return http.StatusServiceUnavailable
}

func (r AddressableNotFoundByNameErrorConcept) isA(err error) bool {
	if err == db.ErrNotFound {
		return false
	}
	return true
}
