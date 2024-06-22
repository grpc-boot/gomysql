package helper

import "strings"

func Replace(table string, columns Columns, rows []Row) (sql string, args []any) {
	if len(columns) == 0 || len(rows) == 0 {
		return
	}

	var (
		buffer      strings.Builder
		fields      = strings.Join(columns, ",")
		placeHolder = repeatAndJoin("?", ",", len(columns))
	)
	length := 7 + 6 + len(table) + len(fields) + 2 + 7 + 2 + len(placeHolder) + (len(rows)-1)*(len(placeHolder)+3)

	buffer.Grow(length)

	args = make([]any, 0, len(columns)*len(rows))

	buffer.WriteString("REPLACE INTO ")
	buffer.WriteString(table)
	buffer.WriteByte('(')
	buffer.WriteString(fields)
	buffer.WriteByte(')')
	buffer.WriteString(" VALUES")
	buffer.WriteByte('(')
	buffer.WriteString(placeHolder)
	buffer.WriteByte(')')

	args = append(args, rows[0]...)
	for index := 1; index < len(rows); index++ {
		buffer.WriteString(",(")
		buffer.WriteString(placeHolder)
		buffer.WriteByte(')')
		args = append(args, rows[index]...)
	}

	return buffer.String(), args
}
