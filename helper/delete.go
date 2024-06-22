package helper

import (
	"strings"

	"github.com/grpc-boot/gomysql/condition"
)

func Delete(table string, condition condition.Condition) (sql string, args []any) {
	var (
		where  string
		buffer strings.Builder
	)

	where, args = condition.Build()

	length := 12 + len(table) + len(where)
	if where != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(table)
	if where != "" {
		buffer.WriteString(" WHERE ")
		buffer.WriteString(where)
	}

	return buffer.String(), args
}
