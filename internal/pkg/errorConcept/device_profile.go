package errorConcept

import (
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
)

// DeviceProfileErrorConcept represents the accessor for the device-profile-specific error concepts
type DeviceProfileErrorConcept struct {
	DuplicateName              deviceProfileDuplicateName
	NotFound                   deviceProfileNotFound
	InvalidState               deviceProfileInvalidState
	MissingFile                deviceProfileMissingFile
	InvalidState_BadRequest    deviceProfileInvalidState_BadRequest
	ContractInvalid            deviceProfileContractInvalid
	EmptyName                  deviceProfileEmptyName
	ContractInvalid_BadRequest deviceProfileContractInvalid_BadRequest
}

type deviceProfileDuplicateName struct {
	err error
}

func (r deviceProfileDuplicateName) httpError() int {
	return http.StatusConflict
}

func (r deviceProfileDuplicateName) isA(err error) bool {
	panic("this is a default method")
}

type deviceProfileNotFound struct{}

func (r deviceProfileNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r deviceProfileNotFound) isA(err error) bool {
	_, ok := err.(errors.ErrDeviceProfileNotFound)
	return ok
}

type deviceProfileInvalidState struct{}

func (r deviceProfileInvalidState) httpErrorCode() int {
	return http.StatusConflict
}

func (r deviceProfileInvalidState) isA(err error) bool {
	_, ok := err.(errors.ErrDeviceProfileInvalidState)
	return ok
}

type deviceProfileMissingFile struct{}

func (r deviceProfileMissingFile) httpErrorCode() int {
	return http.StatusBadRequest
}

// TODO Custom Error
func (r deviceProfileMissingFile) isA(err error) bool {
	return err == errors.NewErrEmptyFile("YAML")
}

type deviceProfileInvalidState_BadRequest struct{}

func (r deviceProfileInvalidState_BadRequest) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r deviceProfileInvalidState_BadRequest) isA(err error) bool {
	_, ok := err.(errors.ErrDeviceProfileInvalidState)
	return ok
}

type deviceProfileContractInvalid struct{}

func (r deviceProfileContractInvalid) httpErrorCode() int {
	return http.StatusConflict
}

func (r deviceProfileContractInvalid) isA(err error) bool {
	_, ok := err.(models.ErrContractInvalid)
	return ok
}

type deviceProfileEmptyName struct{}

func (r deviceProfileEmptyName) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r deviceProfileEmptyName) isA(err error) bool {
	_, ok := err.(errors.ErrEmptyDeviceProfileName)
	return ok
}

type deviceProfileContractInvalid_BadRequest struct{}

func (r deviceProfileContractInvalid_BadRequest) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r deviceProfileContractInvalid_BadRequest) isA(err error) bool {
	_, ok := err.(models.ErrContractInvalid)
	return ok
}

//type DeviceProfileServiceClientErrorConcept struct {
//	Err error
//}
//type DeviceProfileNotFoundErrorConcept struct{}
//type ValueDescriptorsInUseErrorConcept struct{}
//type DeviceProfileInvalidStateErrorConcept struct{}
//type DeviceProfileBadRequestErrorConcept struct{}
//type DeviceProfileEmptyNameErrorConcept struct{}
//type DeviceProfileContractInvalidErrorConcept struct{}
//type DeviceProfileMissingFileErrorConcept struct{}
//
//func (r DeviceProfileServiceClientErrorConcept) httpErrorCode() int {
//	return r.Err.(types.ErrServiceClient).StatusCode
//}
//
//func (r DeviceProfileServiceClientErrorConcept) isA(err error) bool {
//	_, ok := err.(types.ErrServiceClient)
//	return ok
//}
//
//func (r DeviceProfileNotFoundErrorConcept) httpErrorCode() int {
//	return http.StatusNotFound
//}
//
//func (r DeviceProfileNotFoundErrorConcept) isA(err error) bool {
//	_, ok := err.(errors.ErrDeviceProfileNotFound)
//	return ok
//}
//
//func (r ValueDescriptorsInUseErrorConcept) httpErrorCode() int {
//	return http.StatusConflict
//}
//
//func (r ValueDescriptorsInUseErrorConcept) isA(err error) bool {
//	_, ok := err.(dataErrors.ErrValueDescriptorsInUse)
//	return ok
//}
//
//func (r DeviceProfileInvalidStateErrorConcept) httpErrorCode() int {
//	return http.StatusConflict
//}
//
//func (r DeviceProfileInvalidStateErrorConcept) isA(err error) bool {
//	_, ok := err.(errors.ErrDeviceProfileInvalidState)
//	return ok
//}
//
//func (r DeviceProfileBadRequestErrorConcept) httpErrorCode() int {
//	return http.StatusBadRequest
//}
//
//func (r DeviceProfileBadRequestErrorConcept) isA(err error) bool {
//	_, ok := err.(errors.ErrDeviceProfileInvalidState)
//	return ok
//}
//
//func (r DeviceProfileEmptyNameErrorConcept) httpErrorCode() int {
//	return http.StatusBadRequest
//}
//
//func (r DeviceProfileEmptyNameErrorConcept) isA(err error) bool {
//	_, ok := err.(errors.ErrEmptyDeviceProfileName)
//	return ok
//}
//
//func (r DeviceProfileContractInvalidErrorConcept) httpErrorCode() int {
//	return http.StatusConflict
//}
//
//func (r DeviceProfileContractInvalidErrorConcept) isA(err error) bool {
//	_, ok := err.(models.ErrContractInvalid)
//	return ok
//}
//
//func (r DeviceProfileMissingFileErrorConcept) httpErrorCode() int {
//	return http.StatusBadRequest
//}
//
//func (r DeviceProfileMissingFileErrorConcept) isA(err error) bool {
//	return err == errors.NewErrEmptyFile("YAML")
//}
