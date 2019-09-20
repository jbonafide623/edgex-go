package httperror

import (
	dataErrors "github.com/edgexfoundry/edgex-go/internal/core/data/errors"
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/types"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"net/http"
)

type DeviceProfileServiceClientErrorConcept struct {
	Err error
}

func (r DeviceProfileServiceClientErrorConcept) httpErrorCode() int {
	return r.Err.(types.ErrServiceClient).StatusCode
}

func (r DeviceProfileServiceClientErrorConcept) isA(err error) bool {
	_, ok := err.(types.ErrServiceClient)
	return ok
}

type DeviceProfileNotFoundErrorConcept struct{}

func (r DeviceProfileNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DeviceProfileNotFoundErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrDeviceProfileNotFound)
	return ok
}

type ValueDescriptorsInUseErrorConcept struct{}

func (r ValueDescriptorsInUseErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r ValueDescriptorsInUseErrorConcept) isA(err error) bool {
	_, ok := err.(dataErrors.ErrValueDescriptorsInUse)
	return ok
}

type DeviceProfileInvalidStateErrorConcept struct{}

func (r DeviceProfileInvalidStateErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DeviceProfileInvalidStateErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrDeviceProfileInvalidState)
	return ok
}

type DuplicateDeviceProfileErrorConcept struct{}

func (r DuplicateDeviceProfileErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DuplicateDeviceProfileErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrDuplicateName)
	return ok
}

type EmptyDeviceProfileNameErrorConcept struct{}

func (r EmptyDeviceProfileNameErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r EmptyDeviceProfileNameErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrEmptyDeviceProfileName)
	return ok
}

type DeviceProfileContractInvalidErrorConcept struct{}

func (r DeviceProfileContractInvalidErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DeviceProfileContractInvalidErrorConcept) isA(err error) bool {
	_, ok := err.(models.ErrContractInvalid)
	return ok
}
