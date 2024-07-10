package gomysql

type PoolOptions struct {
	Masters []Options `json:"masters" yaml:"masters"`
	Slaves  []Options `json:"slaves" yaml:"slaves"`
}
