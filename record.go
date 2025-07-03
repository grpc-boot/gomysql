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

func (r Record) ToUint16(key string) uint16 {
	value, _ := r[key]
	return convert.String2Uint16(value)
}

func (r Record) ToInt16(key string) int16 {
	value, _ := r[key]
	return convert.String2Int16(value)
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

func (r Record) Clone() Record {
	rd := make(Record, len(r))
	for k, v := range rd {
		rd[k] = v
	}
	return rd
}

func (r Record) ToBytesRecord() BytesRecord {
	if r == nil {
		return nil
	}

	if len(r) == 0 {
		return BytesRecord{}
	}

	br := make(BytesRecord, len(r))
	for key, _ := range r {
		br[key] = r.ToBytes(key)
	}

	return br
}
