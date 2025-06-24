package gomysql

import (
	"sync"
	"time"
)

var (
	DefaultCheckInterval = time.Second * 8
)

var (
	_container sync.Map
)

func PutWithOption(key string, opt PoolOptions) error {
	p, err := NewPool(opt, DefaultCheckInterval)
	if err != nil {
		return err
	}

	Put(key, p)
	return nil
}

func Put(key string, p *Pool) {
	_container.Store(key, p)
}

func Get(key string) (p *Pool) {
	var (
		value, _ = _container.Load(key)
	)

	p, _ = value.(*Pool)
	return
}

func RandDb(key string, t DbType) *Db {
	p := Get(key)

	if p == nil {
		return nil
	}

	return p.Rand(t)
}

func RandExecutor(key string, t DbType) Executor {
	p := Get(key)
	if p == nil {
		return nil
	}

	return p.RandExecutor(t)
}
