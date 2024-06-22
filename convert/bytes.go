package convert

import (
	"unsafe"
)

// Bytes2String 字节切片转换为字符串
func Bytes2String(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

// Bytes2Int64 字节切片转换为int64
func Bytes2Int64(data []byte) int64 {
	return String2Int64(Bytes2String(data))
}

// Bytes2Uint64 字节切片转换为uint64
func Bytes2Uint64(data []byte) uint64 {
	return String2Uint64(Bytes2String(data))
}

// Bytes2Int32 字节切片转换为int32
func Bytes2Int32(data []byte) int32 {
	return String2Int32(Bytes2String(data))
}

// Bytes2Uint32 字节切片转换为uint32
func Bytes2Uint32(data []byte) uint32 {
	return String2Uint32(Bytes2String(data))
}

// Bytes2Int16 字节切片转换为int16
func Bytes2Int16(data []byte) int16 {
	return String2Int16(Bytes2String(data))
}

// Bytes2Uint16 字节切片转换为uint16
func Bytes2Uint16(data []byte) uint16 {
	return String2Uint16(Bytes2String(data))
}

// Bytes2Int8 字节切片转换为int8
func Bytes2Int8(data []byte) int8 {
	return String2Int8(Bytes2String(data))
}

// Bytes2Uint8 字节切片转换为uint8
func Bytes2Uint8(data []byte) uint8 {
	return String2Uint8(Bytes2String(data))
}

// Bytes2Float64 字节切片转换为float64
func Bytes2Float64(data []byte) float64 {
	return String2Float64(Bytes2String(data))
}

// Bytes2Bool 字节切片转换为bool
func Bytes2Bool(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	return String2Bool(Bytes2String(data))
}
