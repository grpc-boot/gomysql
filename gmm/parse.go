package gmm

import (
	"regexp"
	"strings"
)

type ColumnInfo struct {
	Field      string
	GoType     string
	Comment    string
	DefaultVal *string
}

func ParseCreateTable(sqlStr string) (cols []ColumnInfo, primaryKeys []string) {
	lines := strings.Split(sqlStr, "\n")

	var (
		re    = regexp.MustCompile("^\\s*`([^`]+)`\\s+([^\\s]+(?:\\s+unsigned)?)")
		rePri = regexp.MustCompile(`(?i)PRIMARY\s+KEY\s*\(([^)]+)\)`)
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(primaryKeys) == 0 {
			pm := rePri.FindStringSubmatch(line)
			if len(pm) > 1 {
				primaryKeys = strings.Split(pm[1], ",")
				for i := range primaryKeys {
					primaryKeys[i] = strings.Trim(strings.TrimSpace(primaryKeys[i]), "`")
				}
			}
		}

		if !strings.HasPrefix(line, "`") {
			continue
		}

		m := re.FindStringSubmatch(line)
		if len(m) < 3 {
			continue
		}

		field := m[1]
		sqlType := strings.ToLower(m[2])

		// 提取 DEFAULT 值
		var def *string
		if i := strings.Index(strings.ToUpper(line), "DEFAULT"); i != -1 {
			// 取 DEFAULT 后面的部分
			part := line[i+8:]
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "'") {
				end := strings.Index(part[1:], "'")
				if end != -1 {
					val := part[1 : end+1]
					def = &val
				}
			} else {
				// 如 CURRENT_TIMESTAMP
				fields := strings.Fields(part)
				if len(fields) > 0 {
					v := strings.Trim(fields[0], ",")
					def = &v
				}
			}
		}

		// 提取 COMMENT
		var comment string
		if i := strings.Index(strings.ToUpper(line), "COMMENT"); i != -1 {
			part := line[i+7:]
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "'") {
				end := strings.LastIndex(part, "'")
				if end > 1 {
					comment = part[1:end]
				}
			}
		}

		cols = append(cols, ColumnInfo{
			Field:      field,
			GoType:     SQLTypeToGo(sqlType),
			Comment:    comment,
			DefaultVal: def,
		})
	}

	return
}

// SQLTypeToGo SQL类型 → Go类型
func SQLTypeToGo(sqlType string) string {
	t := strings.ToLower(sqlType)
	switch {
	case strings.HasPrefix(t, "tinyint"):
		if strings.Contains(t, "unsigned") {
			return "uint8"
		}
		return "int8"
	case strings.HasPrefix(t, "smallint"):
		if strings.Contains(t, "unsigned") {
			return "uint16"
		}
		return "int16"
	case strings.HasPrefix(t, "mediumint"), strings.HasPrefix(t, "int"):
		if strings.Contains(t, "unsigned") {
			return "uint32"
		}
		return "int32"
	case strings.HasPrefix(t, "bigint"):
		if strings.Contains(t, "unsigned") {
			return "uint64"
		}
		return "int64"
	case strings.HasPrefix(t, "decimal"), strings.HasPrefix(t, "numeric"),
		strings.HasPrefix(t, "float"), strings.HasPrefix(t, "double"):
		return "float64"
	case strings.Contains(t, "char"), strings.Contains(t, "text"),
		strings.Contains(t, "enum"), strings.Contains(t, "set"):
		return "string"
	case strings.HasPrefix(t, "bool"):
		return "bool"
	case strings.HasPrefix(t, "binary"), strings.HasPrefix(t, "blob"):
		return "[]byte"
	default:
		return "string"
	}
}
