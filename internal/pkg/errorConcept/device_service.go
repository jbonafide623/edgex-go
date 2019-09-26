package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/pkg/errors"
)

// DeviceServiceErrorConcept represents the accessor for the device-service-specific error concepts
type DeviceServiceErrorConcept struct {
	EmptyAddressable    deviceServiceEmptyAddressable
	AddressableNotFound deviceServiceAddressableNotFound
	NotUnique           deviceServiceNotUnique
	NotFound            deviceServiceNotFound
	InvalidState        deviceServiceInvalidState
}

// NewDeviceServiceDuplicate instantiates a new deviceServiceDuplicate error concept in effort to handle stateful
// DeviceService duplicate errors
func NewDeviceServiceDuplicate(currentDSId string, newDSId string) deviceServiceDuplicate {
	return deviceServiceDuplicate{CurrentDSId: currentDSId, NewDSId: newDSId}
}

// deviceServiceDuplicate implements ExplicitErrorConceptType
type deviceServiceDuplicate struct {
	CurrentDSId string
	NewDSId     string
}

func (r deviceServiceDuplicate) httpErrorCode() int {
	return http.StatusConflict
}

func (r deviceServiceDuplicate) isA(err error) bool {
	return r.CurrentDSId != r.NewDSId
}

func (r deviceServiceDuplicate) httpError(err error) error {
	return errors.New("Duplicate name for Device Service")
}

func (r deviceServiceDuplicate) logMessage(err error) string {
	return err.Error()
}

// deviceServiceEmptyAddressable implements ExplicitErrorConceptType
type deviceServiceEmptyAddressable struct{}

func (r deviceServiceEmptyAddressable) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r deviceServiceEmptyAddressable) httpError(err error) error {
	return errors.New("Must provide an Addressable for Device Service")
}

func (r deviceServiceEmptyAddressable) isA(err error) bool {
	panic("isA should not be invoked, this is a fallback error concept only")
}

func (r deviceServiceEmptyAddressable) logMessage(err error) string {
	return err.Error()
}

// deviceServiceAddressableNotFound implements ExplicitErrorConceptType
type deviceServiceAddressableNotFound struct{}

func (r deviceServiceAddressableNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r deviceServiceAddressableNotFound) httpError(err error) error {
	return errors.New("Addressable not found by ID or Name")
}

func (r deviceServiceAddressableNotFound) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r deviceServiceAddressableNotFound) logMessage(err error) string {
	return err.Error()
}

// deviceServiceNotUnique implements ExplicitErrorConceptType
type deviceServiceNotUnique struct{}

func (r deviceServiceNotUnique) httpErrorCode() int {
	return http.StatusConflict
}

func (r deviceServiceNotUnique) httpError(err error) error {
	return errors.New("Duplicate name for the device service")
}

func (r deviceServiceNotUnique) isA(err error) bool {
	return err == db.ErrNotUnique
}

func (r deviceServiceNotUnique) logMessage(err error) string {
	return err.Error()
}

// deviceServiceNotFound implements ExplicitErrorConceptType
type deviceServiceNotFound struct{}

func (r deviceServiceNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r deviceServiceNotFound) httpError(err error) error {
	return errors.New("Device service not found")
}

func (r deviceServiceNotFound) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r deviceServiceNotFound) logMessage(err error) string {
	return err.Error()
}

// deviceServiceInvalidState implements ExplicitErrorConceptType
type deviceServiceInvalidState struct{}

func (r deviceServiceInvalidState) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r deviceServiceInvalidState) httpError(err error) error {
	return err
}

func (r deviceServiceInvalidState) isA(err error) bool {
	panic("isA should not be implemented, this is fallback error concept")
}

func (r deviceServiceInvalidState) logMessage(err error) string {
	return err.Error()
}
