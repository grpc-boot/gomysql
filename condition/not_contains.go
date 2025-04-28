package condition

import "strings"

type NotContains struct {
	Field string `json:"field"`
	Words string `json:"words"`
}

func (nc NotContains) Build() (sql string, args []any) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(nc.Field) + 11)

	buffer.WriteString(nc.Field)
	buffer.WriteString(" NOT LIKE ?")

	return buffer.String(), []any{"%" + nc.Words + "%"}
}
