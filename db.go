package gomysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/grpc-boot/gomysql/helper"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	opt  Options
	pool *sql.DB
}

func NewDb(opt Options) (*Db, error) {
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

	p.SetConnMaxLifetime(opt.ConnMaxLifetime())
	p.SetConnMaxIdleTime(opt.ConnMaxIdleTime())

	return &Db{
		opt:  opt,
		pool: p,
	}, nil
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
	return db.FindOneContext(context.Background(), q)
}

func (db *Db) FindOneContext(ctx context.Context, q *helper.Query) (Record, error) {
	records, err := db.FindContext(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records[0], nil
}

func (db *Db) FindOneTimeout(timeout time.Duration, q *helper.Query) (Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return db.FindOneContext(ctx, q)
}

func (db *Db) Find(q *helper.Query) (records []Record, err error) {
	rows, err := Select(db.pool, q)
	return Scan(rows, err)
}

func (db *Db) FindContext(ctx context.Context, q *helper.Query) (records []Record, err error) {
	rows, err := SelectContext(ctx, db.pool, q)
	return Scan(rows, err)
}

func (db *Db) FindTimeout(timeout time.Duration, q *helper.Query) (records []Record, err error) {
	rows, err := SelectTimeout(timeout, db.pool, q)
	return Scan(rows, err)
}
