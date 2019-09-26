package errorConcept

import (
	"errors"
	"net/http"

	metadataErrors "github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
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

// deviceProfileDuplicateName implements ExplicitErrorConceptType
type deviceProfileDuplicateName struct {
	err error
}

func (r deviceProfileDuplicateName) httpErrorCode() int {
	return http.StatusConflict
}

func (r deviceProfileDuplicateName) isA(err error) bool {
	panic("this is a default method")
}

func (r deviceProfileDuplicateName) httpError(err error) error {
	return errors.New("Duplicate name for device profile")
}

func (r deviceProfileDuplicateName) logMessage(err error) string {
	return err.Error()
}

type deviceProfileNotFound struct{}

func (r deviceProfileNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r deviceProfileNotFound) isA(err error) bool {
	_, ok := err.(metadataErrors.ErrDeviceProfileNotFound)
	return ok
}

type deviceProfileInvalidState struct{}

func (r deviceProfileInvalidState) httpErrorCode() int {
	return http.StatusConflict
}

func (r deviceProfileInvalidState) isA(err error) bool {
	_, ok := err.(metadataErrors.ErrDeviceProfileInvalidState)
	return ok
}

// deviceProfileMissingFile implements ExplicitErrorConceptType
type deviceProfileMissingFile struct{}

func (r deviceProfileMissingFile) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r deviceProfileMissingFile) isA(err error) bool {
	return err == http.ErrMissingFile
}

func (r deviceProfileMissingFile) httpError(err error) error {
	return metadataErrors.NewErrEmptyFile("YAML")
}

func (r deviceProfileMissingFile) logMessage(err error) string {
	return err.Error()
}

type deviceProfileInvalidState_BadRequest struct{}

func (r deviceProfileInvalidState_BadRequest) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r deviceProfileInvalidState_BadRequest) isA(err error) bool {
	_, ok := err.(metadataErrors.ErrDeviceProfileInvalidState)
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
	_, ok := err.(metadataErrors.ErrEmptyDeviceProfileName)
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
