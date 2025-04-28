package helper

import "math"

func PageCount(totalCount, pageSize int64) int64 {
	return int64(math.Ceil(float64(totalCount) / float64(pageSize)))
}
