package convert

import (
	"strconv"
	"unsafe"
)

// String2Bytes 字符串转字节切片，注意：转换后不能对字节切片进行修改
func String2Bytes(data string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&data))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// String2Int64 _
func String2Int64(data string) int64 {
	val, _ := strconv.ParseInt(data, 10, 64)
	return val
}

// String2Uint64 _
func String2Uint64(data string) uint64 {
	val, _ := strconv.ParseUint(data, 10, 64)
	return val
}

// String2Int32 _
func String2Int32(data string) int32 {
	val, _ := strconv.ParseInt(data, 10, 32)
	return int32(val)
}

// String2Uint32 _
func String2Uint32(data string) uint32 {
	val, _ := strconv.ParseUint(data, 10, 32)
	return uint32(val)
}

// String2Int16 _
func String2Int16(data string) int16 {
	val, _ := strconv.ParseInt(data, 10, 16)
	return int16(val)
}

// String2Uint16 _
func String2Uint16(data string) uint16 {
	val, _ := strconv.ParseUint(data, 10, 16)
	return uint16(val)
}

// String2Int8 _
func String2Int8(data string) int8 {
	val, _ := strconv.ParseInt(data, 10, 8)
	return int8(val)
}

// String2Uint8 _
func String2Uint8(data string) uint8 {
	val, _ := strconv.ParseUint(data, 10, 8)
	return uint8(val)
}

// String2Float64 _
func String2Float64(data string) float64 {
	val, _ := strconv.ParseFloat(data, 64)
	return val
}

// String2Bool _
func String2Bool(data string) bool {
	if len(data) == 0 {
		return false
	}

	val, _ := strconv.ParseBool(data)
	return val
}
