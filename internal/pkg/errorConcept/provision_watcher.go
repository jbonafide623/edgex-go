package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/pkg/errors"
)

// ProvisionWatcherErrorConcept represents the accessor for provision-watcher-specific error concepts
type ProvisionWatcherErrorConcept struct {
	DeviceServiceNotFound          provisionWatcherDSNotFound
	DeviceServiceNotFound_Conflict provisionWatcherDeviceServiceNotFound_Conflict
	NotFoundById                   provisionWatcherNotFoundById
	NotFoundByName                 provisionWatcherNotFoundByName
	ServiceUnavailable             provisionWatcherServiceUnavailable
	NotUnique                      provisionWatcherNotUnique
	DeviceProfileNotFound          provisionWatcherDeviceProfileNotFound
	DeviceProfileNotFound_Conflict provisionWatcherDeviceProfileNotFound_Conflict
}

// NewProvisionWatcherDuplicateErrorConcept instantiates a new error concept for a given set of ids
func NewProvisionWatcherDuplicateErrorConcept(currentId string, newId string) provisionWatcherDuplicate {
	return provisionWatcherDuplicate{currentId: currentId, newId: newId}
}

// implements ExplicitErrorConceptType
type provisionWatcherDuplicate struct {
	currentId string
	newId     string
}

func (r provisionWatcherDuplicate) httpErrorCode() int {
	return http.StatusConflict
}

func (r provisionWatcherDuplicate) isA(err error) bool {
	return err != db.ErrNotFound && r.currentId != r.newId
}

func (r provisionWatcherDuplicate) httpError(err error) error {
	return errors.New("Duplicate name for the provision watcher")
}

func (r provisionWatcherDuplicate) logMessage(err error) string {
	return err.Error()
}

// provisionWatcherDSNotFound implements ExplicitErrorConceptType
type provisionWatcherDSNotFound struct{}

func (r provisionWatcherDSNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r provisionWatcherDSNotFound) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r provisionWatcherDSNotFound) httpError(err error) error {
	return errors.New("Device service not found")
}

func (r provisionWatcherDSNotFound) logMessage(err error) string {
	return "Device service not found: " + err.Error()
}

// provisionWatcherNotFoundById implements ExplicitErrorConceptType
type provisionWatcherNotFoundById struct{}

func (r provisionWatcherNotFoundById) httpErrorCode() int {
	return http.StatusNotFound
}

func (r provisionWatcherNotFoundById) httpError(err error) error {
	return errors.New("Provision Watcher not found by ID: " + err.Error())
}

func (r provisionWatcherNotFoundById) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r provisionWatcherNotFoundById) logMessage(err error) string {
	return "Provision Watcher not found by ID: " + err.Error()
}

// provisionWatcherNotFoundByName implements ExplicitErrorConceptType
type provisionWatcherNotFoundByName struct{}

func (r provisionWatcherNotFoundByName) httpErrorCode() int {
	return http.StatusNotFound
}

func (r provisionWatcherNotFoundByName) httpError(err error) error {
	return errors.New("Provision Watcher not found: " + err.Error())
}

func (r provisionWatcherNotFoundByName) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r provisionWatcherNotFoundByName) logMessage(err error) string {
	return "Provision Watcher not found: " + err.Error()
}

// provisionWatcherServiceUnavailable implements ExplicitErrorConceptType
type provisionWatcherServiceUnavailable struct{}

func (r provisionWatcherServiceUnavailable) httpErrorCode() int {
	return http.StatusServiceUnavailable
}

func (r provisionWatcherServiceUnavailable) httpError(err error) error {
	return err
}

func (r provisionWatcherServiceUnavailable) isA(err error) bool {
	panic("isA should not be invoke, this is a default error concept")
}

func (r provisionWatcherServiceUnavailable) logMessage(err error) string {
	return "Problem getting provision watcher: " + err.Error()
}

// provisionWatcherDeviceServiceNotFound_Conflict implements ExplicitErrorConceptType
type provisionWatcherDeviceServiceNotFound_Conflict struct{}

func (r provisionWatcherDeviceServiceNotFound_Conflict) httpErrorCode() int {
	return http.StatusConflict
}

func (r provisionWatcherDeviceServiceNotFound_Conflict) httpError(err error) error {
	return errors.New("Device service not found for provision watcher")
}

func (r provisionWatcherDeviceServiceNotFound_Conflict) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r provisionWatcherDeviceServiceNotFound_Conflict) logMessage(err error) string {
	return "Device service not found for provision watcher: " + err.Error()
}

// provisionWatcherNotUnique implements ExplicitErrorConceptType
type provisionWatcherNotUnique struct{}

func (r provisionWatcherNotUnique) httpErrorCode() int {
	return http.StatusConflict
}

func (r provisionWatcherNotUnique) httpError(err error) error {
	return errors.New("Duplicate name for the provision watcher")
}

func (r provisionWatcherNotUnique) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r provisionWatcherNotUnique) logMessage(err error) string {
	return "Duplicate name for the provision watcher: " + err.Error()
}

// provisionWatcherDeviceProfileNotFound implements ExplicitErrorConceptType
type provisionWatcherDeviceProfileNotFound struct{}

func (r provisionWatcherDeviceProfileNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r provisionWatcherDeviceProfileNotFound) httpError(err error) error {
	return errors.New("Device profile not found")
}

func (r provisionWatcherDeviceProfileNotFound) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r provisionWatcherDeviceProfileNotFound) logMessage(err error) string {
	return "Device profile not found: " + err.Error()
}

// provisionWatcherDeviceProfileNotFound_Conflict implements ExplicitErrorConceptType
type provisionWatcherDeviceProfileNotFound_Conflict struct{}

func (r provisionWatcherDeviceProfileNotFound_Conflict) httpErrorCode() int {
	return http.StatusConflict
}

func (r provisionWatcherDeviceProfileNotFound_Conflict) httpError(err error) error {
	return errors.New("Device profile not found for provision watcher")
}

func (r provisionWatcherDeviceProfileNotFound_Conflict) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r provisionWatcherDeviceProfileNotFound_Conflict) logMessage(err error) string {
	return "Device profile not found for provision watcher: " + err.Error()
}
