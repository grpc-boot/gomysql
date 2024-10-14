package gomysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"golang.org/x/exp/rand"
	"sync/atomic"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

const (
	TypeMaster DbType = 0
	TypeSlave  DbType = 1
)

type DbType uint8

type Pool struct {
	masters       []*Db
	slaves        []*Db
	activeMasters atomic.Value
	activeSlaves  atomic.Value
	latestMaster  atomic.Int32
	latestSlave   atomic.Int32
	ticker        *time.Ticker
}

func NewPool(opt PoolOptions, checkInterval time.Duration) (*Pool, error) {
	var (
		masters = make([]*Db, len(opt.Masters))
		actM    = make([]int, len(opt.Masters))
	)

	for index, o := range opt.Masters {
		db, err := NewDb(o)
		if err != nil {
			return nil, err
		}
		actM[index] = index
		masters[index] = db
	}

	p := &Pool{
		masters: masters,
		ticker:  time.NewTicker(checkInterval),
	}
	p.activeMasters.Store(actM)

	if len(opt.Slaves) == 0 {
		go p.monitor()
		return p, nil
	}

	var (
		slaves = make([]*Db, len(opt.Slaves))
		actS   = make([]int, len(opt.Slaves))
	)

	for index, o := range opt.Slaves {
		db, err := NewDb(o)
		if err != nil {
			return nil, err
		}

		actS[index] = index
		slaves[index] = db
	}

	p.slaves = slaves
	p.activeSlaves.Store(actS)

	go p.monitor()
	return p, nil
}

func (p *Pool) monitor() {
	for range p.ticker.C {
		if len(p.masters) > 0 {
			p.activeMasters.Store(p.checkActive(p.masters))
		}

		if len(p.slaves) > 0 {
			p.activeSlaves.Store(p.checkActive(p.slaves))
		}
	}
}

func (p *Pool) checkActive(list []*Db) []int {
	if len(list) < 2 {
		return nil
	}

	var (
		activeList = make([]int, 0, len(list))
		ch         = make(chan int, 1)
	)

	for index, db := range list {
		go func(i int, d *Db) {
			if d.Pool().Ping() == nil {
				ch <- i
			} else {
				ch <- -1
			}
		}(index, db)
	}

	for i := 0; i < len(list); i++ {
		if index, _ := <-ch; index > -1 {
			activeList = append(activeList, index)
		}
	}

	return activeList
}

func (p *Pool) exec(dbType DbType, fn func(db *Db) error) (err error) {
	if dbType == TypeSlave && len(p.slaves) > 0 {
		if len(p.slaves) == 1 {
			return fn(p.slaves[0])
		}

		al, _ := p.activeSlaves.Load().([]int)
		if len(al) == 0 {
			index := int(p.latestSlave.Load())
			return fn(p.slaves[index])
		}

		rand.Shuffle(len(al), func(i, j int) {
			al[i], al[j] = al[j], al[i]
		})

		for _, index := range al {
			err = fn(p.slaves[index])
			if err != nil && errors.Is(err, driver.ErrBadConn) {
				continue
			}
			p.latestSlave.Store(int32(index))
			return
		}
		return
	}

	if len(p.masters) == 1 {
		return fn(p.masters[0])
	}

	al, _ := p.activeMasters.Load().([]int)
	if len(al) == 0 {
		index := int(p.latestMaster.Load())
		return fn(p.masters[index])
	}

	rand.Shuffle(len(al), func(i, j int) {
		al[i], al[j] = al[j], al[i]
	})

	for _, index := range al {
		err = fn(p.masters[index])
		if err != nil && errors.Is(err, driver.ErrBadConn) {
			continue
		}
		p.latestMaster.Store(int32(index))
		return
	}
	return
}

func (p *Pool) queryRowContext(ctx context.Context, dbType DbType, query string, args ...any) (row *sql.Row) {
	p.exec(dbType, func(db *Db) error {
		row = db.queryRowContext(ctx, query, args...)
		return row.Err()
	})
	return
}

func (p *Pool) queryContext(ctx context.Context, dbType DbType, query string, args ...any) (rows *sql.Rows, err error) {
	p.exec(dbType, func(db *Db) error {
		rows, err = db.queryContext(ctx, query, args...)
		return err
	})
	return
}

func (p *Pool) execContext(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		result, err = db.execContext(ctx, query, args...)
		return err
	})
	return
}

func (p *Pool) AcquireQuery() *helper.Query {
	return helper.AcquireQuery()
}

func (p *Pool) FindOne(dbType DbType, q *helper.Query) (Record, error) {
	return p.FindOneContext(context.Background(), dbType, q)
}

func (p *Pool) FindOneContext(ctx context.Context, dbType DbType, q *helper.Query) (r Record, err error) {
	p.exec(dbType, func(db *Db) error {
		r, err = db.FindOneContext(ctx, q)
		return err
	})
	return
}

func (p *Pool) FindOneTimeout(timeout time.Duration, dbType DbType, q *helper.Query) (Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return p.FindOneContext(ctx, dbType, q)
}

func (p *Pool) Find(dbType DbType, q *helper.Query) (records []Record, err error) {
	return p.FindContext(context.Background(), dbType, q)
}

func (p *Pool) FindContext(ctx context.Context, dbType DbType, q *helper.Query) (records []Record, err error) {
	p.exec(dbType, func(db *Db) error {
		records, err = db.FindContext(ctx, q)
		return err
	})
	return
}

func (p *Pool) FindTimeout(timeout time.Duration, dbType DbType, q *helper.Query) (records []Record, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return p.FindContext(ctx, dbType, q)
}

func (p *Pool) Insert(table string, columns helper.Columns, rows ...helper.Row) (sql.Result, error) {
	return p.InsertContext(context.Background(), table, columns, rows...)
}

func (p *Pool) InsertContext(ctx context.Context, table string, columns helper.Columns, rows ...helper.Row) (result sql.Result, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		result, err = db.InsertContext(ctx, table, columns, rows...)
		return err
	})
	return
}

func (p *Pool) InsertWithInsertedIdContext(ctx context.Context, table string, columns helper.Columns, row helper.Row) (id int64, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		id, err = db.InsertWithInsertedIdContext(ctx, table, columns, row)
		return err
	})

	return
}

func (p *Pool) Update(table string, setter string, where condition.Condition, setterArgs ...any) (sql.Result, error) {
	return p.UpdateContext(context.Background(), table, setter, where, setterArgs...)
}

func (p *Pool) UpdateContext(ctx context.Context, table string, setter string, where condition.Condition, setterArgs ...any) (result sql.Result, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		result, err = db.UpdateContext(ctx, table, setter, where, setterArgs...)
		return err
	})
	return
}

func (p *Pool) UpdateWithRowsAffectedContext(ctx context.Context, table string, setter string, where condition.Condition, setterArgs ...any) (rows int64, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		rows, err = db.UpdateWithRowsAffectedContext(ctx, table, setter, where, setterArgs...)
		return err
	})
	return
}

func (p *Pool) Delete(table string, where condition.Condition) (sql.Result, error) {
	return p.DeleteContext(context.Background(), table, where)
}

func (p *Pool) DeleteContext(ctx context.Context, table string, where condition.Condition) (result sql.Result, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		result, err = db.DeleteContext(ctx, table, where)
		return err
	})
	return
}

func (p *Pool) DeleteWithRowsAffectedContext(ctx context.Context, table string, where condition.Condition) (rows int64, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		rows, err = db.DeleteWithRowsAffectedContext(ctx, table, where)
		return err
	})
	return
}

func (p *Pool) Begin() (tx *sql.Tx, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		tx, err = db.Begin()
		return err
	})
	return
}

func (p *Pool) BeginTx(ctx context.Context, opts *sql.TxOptions) (tx *sql.Tx, err error) {
	p.exec(TypeMaster, func(db *Db) error {
		tx, err = db.BeginTx(ctx, opts)
		return err
	})
	return
}

func (p *Pool) Close() error {
	p.ticker.Stop()

	for _, db := range p.masters {
		db.Pool().Close()
	}

	if len(p.slaves) > 0 {
		for _, db := range p.slaves {
			db.Pool().Close()
		}
	}

	return nil
}

func (p *Pool) Rand(dbType DbType) *Db {
	if dbType == TypeSlave && len(p.slaves) > 0 {
		if len(p.slaves) == 1 {
			return p.slaves[0]
		}

		al, _ := p.activeSlaves.Load().([]int)
		if len(al) == 0 {
			index := int(p.latestSlave.Load())
			return p.slaves[index]
		}

		return p.slaves[al[rand.Intn(len(al))]]
	}

	if len(p.masters) == 1 {
		return p.masters[0]
	}

	al, _ := p.activeMasters.Load().([]int)
	if len(al) == 0 {
		index := int(p.latestMaster.Load())
		return p.masters[index]
	}

	return p.masters[al[rand.Intn(len(al))]]
}

func (p *Pool) RandExecutor(dbType DbType) Executor {
	db := p.Rand(dbType)
	if db != nil {
		return db.Executor()
	}

	return nil
}
