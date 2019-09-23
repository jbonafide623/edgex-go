package httperror

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"net/http"
)

type ProvisionWatcherDatabaseNotFoundErrorConcept struct{}

func (r ProvisionWatcherDatabaseNotFoundErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r ProvisionWatcherDatabaseNotFoundErrorConcept) isA(err error) bool {
	return db.ErrNotFound == err
}

type DuplicateProvisionWatcherErrorConcept struct {
	CurrentPwId string
	UpdatedPwId string
}

func (r DuplicateProvisionWatcherErrorConcept) httpErrorCode() int {
	return http.StatusConflict
}

func (r DuplicateProvisionWatcherErrorConcept) isA(err error) bool {
	return err != db.ErrNotFound && r.CurrentPwId != r.UpdatedPwId
}
