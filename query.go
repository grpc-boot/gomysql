package gomysql

import (
	"context"
	"database/sql"
	"time"
)

func Query(db Executor, query string, args ...any) (*sql.Rows, error) {
	return QueryContext(context.Background(), db, query, args...)
}

func QueryContext(ctx context.Context, db Executor, query string, args ...any) (*sql.Rows, error) {
	defer func() {
		if logger != nil {
			logger(query, args...)
		}
	}()

	return db.QueryContext(ctx, query, args...)
}

func QueryTimeout(timeout time.Duration, db Executor, query string, args ...any) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return QueryContext(ctx, db, query, args...)
}
