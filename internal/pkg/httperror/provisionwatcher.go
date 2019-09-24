package httperror

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"net/http"
)

type ProvisionWatcherDeviceServiceNotFoundErrorConcept struct{}

func (r ProvisionWatcherDeviceServiceNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r ProvisionWatcherDeviceServiceNotFoundErrorConcept) isA(err error) bool {
	return db.ErrNotFound == err
}

type ProvisionWatcherDuplicateErrorConcept struct {
	CurrentPWId string
	UpdatedPWId string
}

func (r ProvisionWatcherDuplicateErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r ProvisionWatcherDuplicateErrorConcept) isA(err error) bool {
	return err != db.ErrNotFound && r.CurrentPWId != r.UpdatedPWId
}
