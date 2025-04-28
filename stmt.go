package gomysql

import (
	"context"
	"database/sql"
)

func Prepare(db Executor, query string) (stmt *sql.Stmt, err error) {
	return PrepareContext(context.Background(), db, query)
}

func PrepareContext(ctx context.Context, db Executor, query string) (stmt *sql.Stmt, err error) {
	stmt, err = db.PrepareContext(ctx, query)
	WriteLog(err, query)
	return
}
