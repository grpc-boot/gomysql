package gomysql

import (
	"database/sql"
	"errors"
)

func IsNil(err error) bool {
	return err == nil || errors.Is(err, sql.ErrNoRows)
}
