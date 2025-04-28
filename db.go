package gomysql

import (
	"context"
	"database/sql"

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

func (db *Db) Executor() Executor {
	return db.pool
}

func (db *Db) AcquireQuery() *helper.Query {
	return helper.AcquireQuery()
}

func (db *Db) Begin() (tx *sql.Tx, err error) {
	tx, err = db.pool.Begin()
	WriteLog(err, "BEGIN")
	return
}

func (db *Db) BeginTx(ctx context.Context, opts *sql.TxOptions) (tx *sql.Tx, err error) {
	tx, err = db.pool.BeginTx(ctx, opts)
	WriteLog(err, "BEGIN", opts.Isolation, opts.ReadOnly)
	return
}
