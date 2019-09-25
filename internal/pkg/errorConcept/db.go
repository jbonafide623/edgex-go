package errorConcept

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
)

// Database represents a collection of database error concepts
type DatabaseErrorConcept struct {
	NotFound      dbNotFound
	NotUnique     dbNotUnique
	InvalidObject dbInvalidObject
}

type dbNotFound struct{}

func (r dbNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

// isA determines if the err is one in which a resource cannot be found in the database
func (r dbNotFound) isA(err error) bool {
	return err == db.ErrNotFound
}

type dbNotUnique struct{}

func (r dbNotUnique) httpErrorCode() int {
	return http.StatusConflict
}

// isA determines if the error is one in which a resource is not unique
func (r dbNotUnique) isA(err error) bool {
	return err == db.ErrNotUnique
}

type dbInvalidObject struct{}

func (r dbInvalidObject) httpErrorCode() int {
	return http.StatusBadRequest
}

// isA determines if the error is one in which the a resource is invalid
func (r dbInvalidObject) isA(err error) bool {
	return err == db.ErrInvalidObjectId
}
