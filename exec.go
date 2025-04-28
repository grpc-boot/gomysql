package gomysql

import (
	"context"
	"database/sql"
)

func Exec(db Executor, query string, args ...any) (sql.Result, error) {
	return ExecContext(context.Background(), db, query, args...)
}

func ExecContext(ctx context.Context, db Executor, query string, args ...any) (res sql.Result, err error) {
	res, err = db.ExecContext(ctx, query, args...)
	WriteLog(err, query, args...)
	return
}

func ExecWithRowsAffectedContext(ctx context.Context, db Executor, query string, args ...any) (rows int64, err error) {
	res, err := ExecContext(ctx, db, query, args...)
	if err != nil {
		return
	}

	return res.RowsAffected()
}

func ExecWithInsertedIdContext(ctx context.Context, db Executor, query string, args ...any) (insertedId int64, err error) {
	res, err := ExecContext(ctx, db, query, args...)
	if err != nil {
		return
	}

	return res.LastInsertId()
}
