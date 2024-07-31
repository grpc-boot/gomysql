package gomysql

import "github.com/grpc-boot/gomysql/convert"

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

func (br BytesRecord) Clone() BytesRecord {
	r := make(BytesRecord, len(br))
	for k, v := range br {
		r[k] = v
	}
	return r
}

func (br BytesRecord) ToRecord() Record {
	if br == nil {
		return nil
	}

	if len(br) == 0 {
		return Record{}
	}

	r := make(Record, len(br))
	for key, _ := range br {
		r[key] = br.String(key)
	}

	return r
}
