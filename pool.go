package gomysql

import (
	"golang.org/x/exp/rand"
	"sync/atomic"
	"time"

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

func (p *Pool) AcquireQuery() *helper.Query {
	return helper.AcquireQuery()
}

func (p *Pool) Close() error {
	p.ticker.Stop()

	for _, db := range p.masters {
		_ = db.Pool().Close()
	}

	if len(p.slaves) > 0 {
		for _, db := range p.slaves {
			_ = db.Pool().Close()
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
