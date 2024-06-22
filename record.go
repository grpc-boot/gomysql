package gomysql

type Record map[string]string

func (r Record) Exists(key string) bool {
	_, exists := r[key]
	return exists
}
