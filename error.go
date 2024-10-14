package gomysql

import (
	"database/sql"
	"errors"
)

func DealNotRowsError(err error) error {
	if err == nil {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return err
}
