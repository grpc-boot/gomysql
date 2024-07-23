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

func (r Record) ToBool(key string) bool {
	value, _ := r[key]
	return convert.String2Bool(value)
}

func (r Record) ToInt(key string) int {
	value, _ := r[key]
	return int(convert.String2Int64(value))
}

func (r Record) ToUint32(key string) uint32 {
	value, _ := r[key]
	return convert.String2Uint32(value)
}

func (r Record) ToInt32(key string) int32 {
	value, _ := r[key]
	return convert.String2Int32(value)
}

func (r Record) ToUint64(key string) uint64 {
	value, _ := r[key]
	return convert.String2Uint64(value)
}

func (r Record) ToInt64(key string) int64 {
	value, _ := r[key]
	return convert.String2Int64(value)
}

func (r Record) ToUint8(key string) uint8 {
	value, _ := r[key]
	return convert.String2Uint8(value)
}

func (r Record) ToInt8(key string) int8 {
	value, _ := r[key]
	return convert.String2Int8(value)
}

func (r Record) ToFloat64(key string) float64 {
	value, _ := r[key]
	return convert.String2Float64(value)
}

func (r Record) ToBytes(key string) []byte {
	value, _ := r[key]
	return convert.String2Bytes(value)
}

type BytesRecord map[string][]byte

func (br BytesRecord) Exists(key string) bool {
	_, exists := br[key]
	return exists
}

func (br BytesRecord) String(key string) string {
	value, _ := br[key]
	return convert.Bytes2String(value)
}

func (br BytesRecord) ToBool(key string) bool {
	value, _ := br[key]
	return convert.Bytes2Bool(value)
}

func (br BytesRecord) ToInt(key string) int {
	value, _ := br[key]
	return int(convert.Bytes2Int64(value))
}

func (br BytesRecord) ToUint32(key string) uint32 {
	value, _ := br[key]
	return convert.Bytes2Uint32(value)
}

func (br BytesRecord) ToInt32(key string) int32 {
	value, _ := br[key]
	return convert.Bytes2Int32(value)
}

func (br BytesRecord) ToUint64(key string) uint64 {
	value, _ := br[key]
	return convert.Bytes2Uint64(value)
}

func (br BytesRecord) ToInt64(key string) int64 {
	value, _ := br[key]
	return convert.Bytes2Int64(value)
}

func (br BytesRecord) ToUint8(key string) uint8 {
	value, _ := br[key]
	return convert.Bytes2Uint8(value)
}

func (br BytesRecord) ToInt8(key string) int8 {
	value, _ := br[key]
	return convert.Bytes2Int8(value)
}

func (br BytesRecord) ToFloat64(key string) float64 {
	value, _ := br[key]
	return convert.Bytes2Float64(value)
}

func (br BytesRecord) Bytes(key string) []byte {
	value, _ := br[key]
	return value
}
