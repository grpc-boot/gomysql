package gomysql

import "github.com/grpc-boot/gomysql/convert"

type Record map[string]string

func (r Record) Exists(key string) bool {
	_, exists := r[key]
	return exists
}

func (r Record) String(key string) string {
	value, _ := r[key]
	return value
}

func (r Record) ToUint32(key string) uint32 {
	value, _ := r[key]
	return convert.String2Uint32(value)
}

func (r Record) ToInt64(key string) int64 {
	value, _ := r[key]
	return convert.String2Int64(value)
}

func (r Record) ToUint8(key string) uint8 {
	value, _ := r[key]
	return convert.String2Uint8(value)
}
