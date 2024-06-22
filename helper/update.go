package helper

import (
	"strings"

	"github.com/grpc-boot/gomysql/condition"
)

func Update(table, setters string, condition condition.Condition) (sql string, args []any) {
	var (
		where  string
		buffer strings.Builder
	)

	where, args = condition.Build()

	length := 7 + len(table) + 5 + len(setters) + len(where)
	if where != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString("UPDATE ")
	buffer.WriteString(table)
	buffer.WriteString(" SET ")
	buffer.WriteString(setters)
	if where != "" {
		buffer.WriteString(" WHERE ")
		buffer.WriteString(where)
	}

	return buffer.String(), args
}
