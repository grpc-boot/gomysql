package gmm

import (
	_ "embed"
	"strings"
)

//go:embed model.tpl
var modelTemplate []byte

//go:embed crud.tpl
var crudTemplate []byte

var (
	ModelTemplate = func() []byte {
		t := make([]byte, len(modelTemplate))
		copy(t, modelTemplate)
		return t
	}

	CrudTemplate = func() []byte {
		t := make([]byte, len(crudTemplate))
		copy(t, crudTemplate)
		return t
	}
)

// 转为大驼峰命名
func bigCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(string(parts[i][0])) + strings.ToLower(parts[i][1:])
		}
	}
	return strings.Join(parts, "")
}

// 转为小驼峰命名
func smallCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			if i == 0 {
				parts[i] = strings.ToLower(parts[i])
				continue
			}

			parts[i] = strings.ToUpper(string(parts[i][0])) + strings.ToLower(parts[i][1:])
		}
	}
	return strings.Join(parts, "")
}
