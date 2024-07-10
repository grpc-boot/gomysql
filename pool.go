package gomysql

import "golang.org/x/exp/rand"

const (
	TypeMaster DbType = 0
	TypeSlave  DbType = 1
)

type DbType uint8

type Pool struct {
	masters []*Db
	slaves  []*Db
}

func NewPool(opt PoolOptions) (*Pool, error) {
	masters := make([]*Db, len(opt.Masters))

	for index, o := range opt.Masters {
		db, err := NewDb(o)
		if err != nil {
			return nil, err
		}

		masters[index] = db
	}

	if len(opt.Slaves) == 0 {
		return &Pool{
			masters: masters,
			slaves:  masters,
		}, nil
	}

	slaves := make([]*Db, len(opt.Slaves))
	for index, o := range opt.Slaves {
		db, err := NewDb(o)
		if err != nil {
			return nil, err
		}

		slaves[index] = db
	}

	return &Pool{
		masters: masters,
		slaves:  slaves,
	}, nil
}

func (p *Pool) Random(dbType DbType) *Db {
	if dbType == TypeMaster {
		return p.masters[rand.Intn(len(p.masters))]
	}
	return p.slaves[rand.Intn(len(p.slaves))]
}

func (p *Pool) Index(dbType DbType, index int) (*Db, error) {
	if dbType == TypeMaster {
		if index >= len(p.masters) {
			return nil, ErrIndexOutRange
		}

		return p.masters[index], nil
	}

	if index >= len(p.slaves) {
		return nil, ErrIndexOutRange
	}
	return p.slaves[index], nil
}
