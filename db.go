package gomysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	opt  Options
	pool *sql.DB
}

func NewDb(opt Options) (*Db, error) {
	opt.format()

	p, err := sql.Open("mysql", opt.Dsn())
	if err != nil {
		return nil, err
	}

	if opt.MaxOpenConns > 0 {
		p.SetMaxOpenConns(opt.MaxOpenConns)
	}

	if opt.MaxIdleConns > 0 {
		p.SetMaxIdleConns(opt.MaxIdleConns)
	}

	if opt.ConnMaxLifetimeSecond > 0 {
		p.SetConnMaxLifetime(opt.ConnMaxLifetime())
	}

	if opt.ConnMaxIdleTimeSecond > 0 {
		p.SetConnMaxIdleTime(opt.ConnMaxIdleTime())
	}

	return &Db{
		opt:  opt,
		pool: p,
	}, nil
}

func (db *Db) Options() Options {
	return db.opt
}

func (db *Db) Pool() *sql.DB {
	return db.pool
}

func (db *Db) Exec(query string, args ...any) (sql.Result, error) {
	return Exec(db.pool, query, args...)
}

func (db *Db) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return ExecContext(ctx, db.pool, query, args...)
}

func (db *Db) ExecTimeout(timeout time.Duration, query string, args ...any) (sql.Result, error) {
	return ExecTimeout(timeout, db.pool, query, args...)
}

func (db *Db) Query(query string, args ...any) (*sql.Rows, error) {
	return Query(db.pool, query, args...)
}

func (db *Db) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return QueryContext(ctx, db.pool, query, args...)
}

func (db *Db) QueryTimeout(timeout time.Duration, query string, args ...any) (*sql.Rows, error) {
	return QueryTimeout(timeout, db.pool, query, args...)
}

func (db *Db) AcquireQuery() *helper.Query {
	return helper.AcquireQuery()
}

func (db *Db) FindOne(q *helper.Query) (Record, error) {
	return FindOne(db.pool, q)
}

func (db *Db) FindOneContext(ctx context.Context, q *helper.Query) (Record, error) {
	return FindOneContext(ctx, db.pool, q)
}

func (db *Db) FindOneTimeout(timeout time.Duration, q *helper.Query) (Record, error) {
	return FindOneTimeout(timeout, db.pool, q)
}

func (db *Db) Find(q *helper.Query) (records []Record, err error) {
	return Find(db.pool, q)
}

func (db *Db) FindContext(ctx context.Context, q *helper.Query) (records []Record, err error) {
	return FindContext(ctx, db.pool, q)
}

func (db *Db) FindTimeout(timeout time.Duration, q *helper.Query) (records []Record, err error) {
	return FindTimeout(timeout, db.pool, q)
}

func (db *Db) Insert(table string, columns helper.Columns, rows ...helper.Row) (sql.Result, error) {
	return Insert(db.pool, table, columns, rows...)
}

func (db *Db) InsertContext(ctx context.Context, table string, columns helper.Columns, rows ...helper.Row) (sql.Result, error) {
	return InsertContext(ctx, db.pool, table, columns, rows...)
}

func (db *Db) InsertTimeout(timeout time.Duration, table string, columns helper.Columns, rows ...helper.Row) (sql.Result, error) {
	return InsertTimeout(timeout, db.pool, table, columns, rows...)
}

func (db *Db) Update(table string, setter string, where condition.Condition, setterArgs ...any) (sql.Result, error) {
	return Update(db.pool, table, setter, where, setterArgs...)
}

func (db *Db) UpdateContext(ctx context.Context, table string, setter string, where condition.Condition, setterArgs ...any) (sql.Result, error) {
	return UpdateContext(ctx, db.pool, table, setter, where, setterArgs...)
}

func (db *Db) UpdateTimeout(timeout time.Duration, table string, setter string, where condition.Condition, setterArgs ...any) (sql.Result, error) {
	return UpdateTimeout(timeout, db.pool, table, setter, where, setterArgs...)
}

func (db *Db) Delete(table string, where condition.Condition) (sql.Result, error) {
	return Delete(db.pool, table, where)
}

func (db *Db) DeleteContext(ctx context.Context, table string, where condition.Condition) (sql.Result, error) {
	return DeleteContext(ctx, db.pool, table, where)
}

func (db *Db) DeleteTimeout(timeout time.Duration, table string, where condition.Condition) (sql.Result, error) {
	return DeleteTimeout(timeout, db.pool, table, where)
}

func (db *Db) Begin() (*sql.Tx, error) {
	return db.pool.Begin()
}

func (db *Db) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.pool.BeginTx(ctx, opts)
}
