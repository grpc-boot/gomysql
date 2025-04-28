package condition

import "strings"

type EndWith struct {
	Field string `json:"field"`
	Words string `json:"words"`
}

func (b EndWith) Build() (sql string, args []any) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(b.Field) + 7)

	buffer.WriteString(b.Field)
	buffer.WriteString(" LIKE ?")

	return buffer.String(), []any{"%" + b.Words}
}
