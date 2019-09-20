package httperror

import (
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"net/http"
)

type DuplicateDeviceErrorConcept struct{}

func (r DuplicateDeviceErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DuplicateDeviceErrorConcept) isA(err error) bool {
	if _, ok := err.(errors.ErrDuplicateName); ok {
		return true
	}

	return false
}

type DeviceNotFoundErrorConcept struct{}

func (r DeviceNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DeviceNotFoundErrorConcept) isA(err error) bool {
	if _, ok := err.(errors.ErrItemNotFound); ok {
		return true
	}
	return false
}

type DeviceDatabaseNotFoundErrorConcept struct{}

func (r DeviceDatabaseNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DeviceDatabaseNotFoundErrorConcept) isA(err error) bool {
	return err == db.ErrNotFound
}

type DeviceDatabaseInvalidObjectErrorConcept struct{}

func (r DeviceDatabaseInvalidObjectErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r DeviceDatabaseInvalidObjectErrorConcept) isA(err error) bool {
	return err == db.ErrInvalidObjectId
}

type DeviceContractInvalidErrorConcept struct{}

func (r DeviceContractInvalidErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r DeviceContractInvalidErrorConcept) isA(err error) bool {
	_, ok := err.(models.ErrContractInvalid)
	return ok
}
