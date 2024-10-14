package gomysql

import (
	"context"
	"database/sql"

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
