package gomysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func Update(db Executor, table, setters string, condition condition.Condition, setterArgs ...any) (sql.Result, error) {
	query, args := helper.Update(table, setters, condition, setterArgs...)
	return Exec(db, query, args...)
}

func UpdateContext(ctx context.Context, db Executor, table, setters string, condition condition.Condition, setterArgs ...any) (sql.Result, error) {
	query, args := helper.Update(table, setters, condition, setterArgs...)
	return ExecContext(ctx, db, query, args...)
}

func UpdateTimeout(timeout time.Duration, db Executor, table, setters string, condition condition.Condition, setterArgs ...any) (sql.Result, error) {
	query, args := helper.Update(table, setters, condition, setterArgs...)
	return ExecTimeout(timeout, db, query, args...)
}
