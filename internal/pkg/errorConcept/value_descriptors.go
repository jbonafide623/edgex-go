package errorConcept

import (
	"github.com/edgexfoundry/edgex-go/internal/core/data/errors"
	"net/http"
)

// ValueDescriptorsErrorConcept represents the accessor for the value-descriptor-specific error concepts
type ValueDescriptorsErrorConcept struct {
	InUse valueDescriptorsInUse
}

type valueDescriptorsInUse struct{}

func (r valueDescriptorsInUse) httpErrorCode() int {
	return http.StatusConflict
}

func (r valueDescriptorsInUse) isA(err error) bool {
	_, ok := err.(errors.ErrValueDescriptorsInUse)
	return ok
}
