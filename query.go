package gomysql

import (
	"context"
	"database/sql"
)

func Query(db Executor, query string, args ...any) (*sql.Rows, error) {
	return QueryContext(context.Background(), db, query, args...)
}

func QueryContext(ctx context.Context, db Executor, query string, args ...any) (rows *sql.Rows, err error) {
	rows, err = db.QueryContext(ctx, query, args...)
	WriteLog(err, query, args...)
	return rows, err
}

func QueryRow(db Executor, query string, args ...any) *sql.Row {
	return QueryRowContext(context.Background(), db, query, args...)
}

func QueryRowContext(ctx context.Context, db Executor, query string, args ...any) (row *sql.Row) {
	var err error
	row = db.QueryRowContext(ctx, query, args...)
	if row != nil {
		err = row.Err()
	}

	WriteLog(err, query, args...)
	return
}
