package gomysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/grpc-boot/gomysql/helper"
)

func Select(db Executor, q *helper.Query) (*sql.Rows, error) {
	query, args := q.Sql()
	return Query(db, query, args...)
}

func SelectContext(ctx context.Context, db Executor, q *helper.Query) (*sql.Rows, error) {
	query, args := q.Sql()
	return QueryContext(ctx, db, query, args...)
}

func SelectTimeout(timeout time.Duration, db Executor, q *helper.Query) (*sql.Rows, error) {
	query, args := q.Sql()
	return QueryTimeout(timeout, db, query, args...)
}
