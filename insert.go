package gomysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/grpc-boot/gomysql/helper"
)

func Insert(db Executor, table string, columns helper.Columns, rows ...helper.Row) (sql.Result, error) {
	query, args := helper.Insert(table, columns, rows, false)
	return Exec(db, query, args...)
}

func InsertContext(ctx context.Context, db Executor, table string, columns helper.Columns, rows ...helper.Row) (sql.Result, error) {
	query, args := helper.Insert(table, columns, rows, false)
	return ExecContext(ctx, db, query, args...)
}

func InsertWithRowsAffectedTimeout(timeout time.Duration, db Executor, table string, columns helper.Columns, rows ...helper.Row) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	query, args := helper.Insert(table, columns, rows, false)
	return ExecWithRowsAffectedContext(ctx, db, query, args...)
}

func InsertWithRowsAffectedContext(ctx context.Context, db Executor, table string, columns helper.Columns, rows ...helper.Row) (rowsAffected int64, err error) {
	query, args := helper.Insert(table, columns, rows, false)
	return ExecWithRowsAffectedContext(ctx, db, query, args...)
}

func InsertWithInsertedIdTimeout(timeout time.Duration, db Executor, table string, columns helper.Columns, row helper.Row) (id int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return InsertWithInsertedIdContext(ctx, db, table, columns, row)
}

func InsertWithInsertedIdContext(ctx context.Context, db Executor, table string, columns helper.Columns, row helper.Row) (id int64, err error) {
	query, args := helper.Insert(table, columns, []helper.Row{row}, false)
	return ExecWithInsertedIdContext(ctx, db, query, args...)
}
