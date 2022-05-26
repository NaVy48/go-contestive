package postgresql

import (
	"contestive/entity"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

const (
	pqCodeCaseNotFound        = "20000"
	pqCodeConstraintViolation = "23505"
)

func mapPqErr(err *pq.Error) error {
	switch err.Code {
	default:
		return err
	case pqCodeCaseNotFound:
		return entity.ErrNotFound(err)
	case pqCodeConstraintViolation:
		return entity.ErrDataConflict(err)
	}
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return entity.ErrNotFound(err)
	}

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	if !ok {
		return err
	}

	return mapPqErr(pqErr)
}
