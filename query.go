package gomysql

import (
	"context"
	"database/sql"
	"time"
)

func Query(db Executor, query string, args ...any) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func QueryContext(ctx context.Context, db Executor, query string, args ...any) (*sql.Rows, error) {
	return db.QueryContext(ctx, query, args...)
}

func QueryTimeout(timeout time.Duration, db Executor, query string, args ...any) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return QueryContext(ctx, db, query, args...)
}
