package gomysql

type Tx interface {
	Commit() error
	Rollback() error
}
