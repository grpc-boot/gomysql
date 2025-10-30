package gomysql

import (
	"database/sql"
	"errors"
)

var (
	ErrModelFileExists = errors.New("model file exists")
	ErrCrudFileExists  = errors.New("crud file exists")
)

func IsNil(err error) bool {
	return err == nil || errors.Is(err, sql.ErrNoRows)
}
