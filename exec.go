package gomysql

import (
	"context"
	"database/sql"
	"time"
)

func Exec(db Executor, query string, args ...any) (sql.Result, error) {
	return db.Exec(query, args...)
}

func ExecContext(ctx context.Context, db Executor, query string, args ...any) (sql.Result, error) {
	return db.ExecContext(ctx, query, args...)
}

func ExecTimeout(timeout time.Duration, db Executor, query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ExecContext(ctx, db, query, args...)
}
