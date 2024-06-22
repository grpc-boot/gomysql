package gomysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/grpc-boot/gomysql/helper"
)

func Insert(db Executor, table string, columns helper.Columns, rows []helper.Row) (sql.Result, error) {
	query, args := helper.Insert(table, columns, rows, false)
	return Exec(db, query, args)
}

func InsertContext(ctx context.Context, db Executor, table string, columns helper.Columns, rows []helper.Row) (sql.Result, error) {
	query, args := helper.Insert(table, columns, rows, false)
	return ExecContext(ctx, db, query, args)
}

func InsertTimeout(timeout time.Duration, db Executor, table string, columns helper.Columns, rows []helper.Row) (sql.Result, error) {
	query, args := helper.Insert(table, columns, rows, false)
	return ExecTimeout(timeout, db, query, args)
}
