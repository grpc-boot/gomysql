package gomysql

import (
	"context"
	"database/sql"
)

func Query(db Executor, query string, args ...any) (*sql.Rows, error) {
	return QueryContext(context.Background(), db, query, args...)
}

func QueryContext(ctx context.Context, db Executor, query string, args ...any) (rows *sql.Rows, err error) {
	defer func() {
		if err != nil && errorLog != nil {
			errorLog(err, query, args...)
		} else if logger != nil {
			logger(query, args...)
		}
	}()

	rows, err = db.QueryContext(ctx, query, args...)
	return rows, err
}

func QueryRow(db Executor, query string, args ...any) *sql.Row {
	return QueryRowContext(context.Background(), db, query, args...)
}

func QueryRowContext(ctx context.Context, db Executor, query string, args ...any) (row *sql.Row) {
	var err error
	defer func() {
		if err != nil && errorLog != nil {
			errorLog(err, query, args...)
		} else if logger != nil {
			logger(query, args...)
		}
	}()

	row = db.QueryRowContext(ctx, query, args...)
	if row != nil {
		err = row.Err()
	}

	return
}
