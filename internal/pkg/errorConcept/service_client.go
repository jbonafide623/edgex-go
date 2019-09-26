package errorConcept

import "github.com/edgexfoundry/go-mod-core-contracts/clients/types"

// TODO Is there a better way?
// NewServiceClientHttpError represents the accessor for the service-client-specific error concepts
func NewServiceClientHttpError(err error) *serviceClientHttpError {
	return &serviceClientHttpError{Err: err}
}

type serviceClientHttpError struct {
	Err error
}

func (r *serviceClientHttpError) httpErrorCode() int {
	return r.Err.(types.ErrServiceClient).StatusCode
}

func (r *serviceClientHttpError) isA(err error) bool {
	_, ok := err.(types.ErrServiceClient)
	return ok
}
