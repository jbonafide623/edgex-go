package httperror

import (
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"net/http"
)

type DeviceServiceDatabaseNotFoundErrorConcept struct{}

func (r DeviceServiceDatabaseNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusServiceUnavailable
}

func (r DeviceServiceDatabaseNotFoundErrorConcept) isA(err error) bool {
	return err == db.ErrNotFound
}

type DuplicateDeviceServiceErrorConcept struct {
	CheckDsId string
	ToId      string
}

func (r DuplicateDeviceServiceErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DuplicateDeviceServiceErrorConcept) isA(err error) bool {
	return r.CheckDsId != r.ToId
}

type DeviceServiceNotFoundErrorConcept struct{}

func (r DeviceServiceNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DeviceServiceNotFoundErrorConcept) isA(err error) bool {
	_, ok := err.(errors.ErrItemNotFound)
	return ok
}
