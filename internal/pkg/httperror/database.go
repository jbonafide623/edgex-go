package httperror

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"net/http"
)

type DatabaseNotUniqueErrorConcept struct{}

func (r DatabaseNotUniqueErrorConcept) httpErrorCode() int {
	return http.StatusNotFound
}

func (r DatabaseNotUniqueErrorConcept) isA(err error) bool {
	return err == db.ErrNotUnique
}
