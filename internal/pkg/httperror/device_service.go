package httperror

import (
	"net/http"
)

type DeviceServiceDuplicateErrorConcept struct {
	OriginalDeviceServiceId string
	UpdatedDeviceServiceId  string
}

func (r DeviceServiceDuplicateErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DeviceServiceDuplicateErrorConcept) isA(err error) bool {
	return r.OriginalDeviceServiceId != r.UpdatedDeviceServiceId
}
