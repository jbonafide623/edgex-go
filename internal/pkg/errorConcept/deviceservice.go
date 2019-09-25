package errorConcept

import (
	"net/http"
)

type DeviceServiceDuplicateErrorConcept struct {
	CurrentDSId string
	UpdatedDSId string
}

func (r DeviceServiceDuplicateErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DeviceServiceDuplicateErrorConcept) isA(err error) bool {
	return r.CurrentDSId != r.UpdatedDSId
}
