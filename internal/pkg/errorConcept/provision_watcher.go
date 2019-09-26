package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
)

// ProvisionWatcherErrorConcept represents the accessor for provision-watcher-specific error concepts
type ProvisionWatcherErrorConcept struct {
	DeviceServiceNotFound provisionWatcherDSNotFound
}

type provisionWatcherDSNotFound struct{}

func (r provisionWatcherDSNotFound) httpErrorCode() int {
	return http.StatusConflict
}

func (r provisionWatcherDSNotFound) isA(err error) bool {
	return err == db.ErrNotFound
}

// NewProvisionWatcherDuplicateErrorConcept instantiates a new error concept for a given set of ids
func NewProvisionWatcherDuplicateErrorConcept(currentId string, newId string) provisionWatcherDuplicate {
	return provisionWatcherDuplicate{currentId: currentId, newId: newId}
}

type provisionWatcherDuplicate struct {
	currentId string
	newId     string
}

func (r provisionWatcherDuplicate) httpErrorCode() int {
	return http.StatusConflict
}

func (r provisionWatcherDuplicate) isA(err error) bool {
	return err != db.ErrNotFound && r.currentId != r.newId
}
