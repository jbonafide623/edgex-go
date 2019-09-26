package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/pkg/errors"
)

// DeviceReportErrorConcept represents the accessor for the device-report-specific error concepts
type DeviceReportErrorConcept struct {
	LimitExceeded deviceReportLimitExceeded
	NotFound      deviceReportNotUnique
	NotUnique     deviceReportNotUnique
}

// deviceReportLimitExceeded implements ExplicitErrorConceptType
type deviceReportLimitExceeded struct{}

func (r deviceReportLimitExceeded) httpErrorCode() int {
	return http.StatusRequestEntityTooLarge
}

func (r deviceReportLimitExceeded) httpError(err error) error {
	return errors.New("Max limit exceeded")
}

func (r deviceReportLimitExceeded) isA(err error) bool {
	panic("isA should not be invoked since it is a default error concept")
}

func (r deviceReportLimitExceeded) logMessage(err error) string {
	return err.Error()
}

// deviceReportDeviceNotFound implements ExplicitErrorConceptType
type deviceReportDeviceNotFound struct{}

func (r deviceReportDeviceNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r deviceReportDeviceNotFound) httpError(err error) error {
	return errors.New("Device referenced by Device Report doesn't exist")
}

func (r deviceReportDeviceNotFound) isA(err error) bool {
	return err == db.ErrNotFound
}

func (r deviceReportDeviceNotFound) logMessage(err error) string {
	return err.Error()
}

// deviceReportNotUnique implements ExplicitErrorConceptType
type deviceReportNotUnique struct{}

func (r deviceReportNotUnique) httpErrorCode() int {
	return http.StatusConflict
}

func (r deviceReportNotUnique) httpError(err error) error {
	return errors.New("Duplicate Name for the device report")
}

func (r deviceReportNotUnique) isA(err error) bool {
	return err == db.ErrNotUnique
}

func (r deviceReportNotUnique) logMessage(err error) string {
	return err.Error()
}
