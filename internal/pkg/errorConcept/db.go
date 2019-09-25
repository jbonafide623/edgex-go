package errorConcept

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"net/http"
)

type Database struct {
	NotFound  dbNotFound
	NotUnique dbNotUnique
}

type dbNotFound struct{}

func (r dbNotFound) httpErrorCode() int {
	return http.StatusNotFound
}

func (r dbNotFound) isA(err error) bool {
	return err == db.ErrNotFound
}

type dbNotUnique struct{}

func (r dbNotUnique) httpErrorCode() int {
	return http.StatusConflict
}

func (r dbNotUnique) isA(err error) bool {
	return err == db.ErrNotUnique
}
