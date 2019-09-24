package httperror

import (
	"net/http"

	dataErrors "github.com/edgexfoundry/edgex-go/internal/core/data/errors"
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/types"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type DeviceProfileServiceClientErrorConcept struct {
	Err error
}
type DeviceProfileNotFoundErrorConcept struct{}
type ValueDescriptorsInUseErrorConcept struct{}
type DeviceProfileInvalidStateErrorConcept struct{}
type DeviceProfileBadRequestErrorConcept struct{}
type DeviceProfileEmptyNameErrorConcept struct{}
type DeviceProfileContractInvalidErrorConcept struct{}
type DeviceProfileMissingFileErrorConcept struct{}

func (r DeviceProfileServiceClientErrorConcept) httpErrorCode() int {
	return r.Err.(types.ErrServiceClient).StatusCode
}

func (r DeviceProfileServiceClientErrorConcept) isA(err error) bool {
	_, ok := err.(types.ErrServiceClient)
	return ok
}

func (r DeviceProfileNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DeviceProfileNotFoundErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrDeviceProfileNotFound)
	return ok
}

func (r ValueDescriptorsInUseErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r ValueDescriptorsInUseErrorConcept) isA(err error) bool {
	_, ok := err.(dataErrors.ErrValueDescriptorsInUse)
	return ok
}

func (r DeviceProfileInvalidStateErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DeviceProfileInvalidStateErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrDeviceProfileInvalidState)
	return ok
}

func (r DeviceProfileBadRequestErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r DeviceProfileBadRequestErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrDeviceProfileInvalidState)
	return ok
}

func (r DeviceProfileEmptyNameErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r DeviceProfileEmptyNameErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrEmptyDeviceProfileName)
	return ok
}

func (r DeviceProfileContractInvalidErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DeviceProfileContractInvalidErrorConcept) isA(err error) bool {
	_, ok := err.(models.ErrContractInvalid)
	return ok
}

func (r DeviceProfileMissingFileErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r DeviceProfileMissingFileErrorConcept) isA(err error) bool {
	return err == errors.NewErrEmptyFile("YAML")
}
