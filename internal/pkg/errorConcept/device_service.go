package errorConcept

import (
	"net/http"
)

// TODO Is there a better way?
// NewDeviceServiceDuplicate instantiates a new deviceServiceDuplicate error concept in effort to handle stateful
// DeviceService duplicate errors
func NewDeviceServiceDuplicate(currentDSId string, newDSId string) deviceServiceDuplicate {
	return deviceServiceDuplicate{CurrentDSId: currentDSId, NewDSId: newDSId}
}

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
