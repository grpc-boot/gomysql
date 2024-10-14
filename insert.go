package gomysql

import (
	"context"
	"database/sql"

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

func InsertWithRowsAffectedContext(ctx context.Context, db Executor, table string, columns helper.Columns, rows ...helper.Row) (rowsAffected int64, err error) {
	query, args := helper.Insert(table, columns, rows, false)
	return ExecWithRowsAffectedContext(ctx, db, query, args...)
}

func InsertWithInsertedIdContext(ctx context.Context, db Executor, table string, columns helper.Columns, row helper.Row) (id int64, err error) {
	query, args := helper.Insert(table, columns, []helper.Row{row}, false)
	return ExecWithInsertedIdContext(ctx, db, query, args...)
}
