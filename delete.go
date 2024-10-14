package gomysql

import (
	"context"
	"database/sql"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func Delete(db Executor, table string, condition condition.Condition) (sql.Result, error) {
	query, args := helper.Delete(table, condition)
	return Exec(db, query, args...)
}

func DeleteContext(ctx context.Context, db Executor, table string, condition condition.Condition) (sql.Result, error) {
	query, args := helper.Delete(table, condition)
	return ExecContext(ctx, db, query, args...)
}

func DeleteWithRowsAffectedContext(ctx context.Context, db Executor, table string, condition condition.Condition) (rows int64, err error) {
	query, args := helper.Delete(table, condition)
	return ExecWithRowsAffectedContext(ctx, db, query, args...)
}
