package httperror

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type DeviceContractInvalidErrorConcept struct{}
type DeviceDatabaseInvalidObjectErrorConcept struct{}

func (r DeviceDatabaseInvalidObjectErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r DeviceDatabaseInvalidObjectErrorConcept) isA(err error) bool {
	return err == db.ErrInvalidObjectId
}

func (r DeviceContractInvalidErrorConcept) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r DeviceContractInvalidErrorConcept) isA(err error) bool {
	_, ok := err.(models.ErrContractInvalid)
	return ok
}
