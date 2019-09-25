package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

// DeviceErrorConcept represents the accessor for the device-specific error concepts
type DeviceErrorConcept struct {
	ContractInvalid contractInvalid
}

type contractInvalid struct{}

func (r contractInvalid) httpErrorCode() int {
	return http.StatusBadRequest
}

func (r contractInvalid) isA(err error) bool {
	_, ok := err.(models.ErrContractInvalid)
	return ok
}
