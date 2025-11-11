package gmm

import (
	"bytes"
	"fmt"
	"strings"
)

// GenerateStruct 生成结构体定义
func GenerateStruct(primaryField, pkg string, table string, cols []ColumnInfo) (model, crud []byte) {
	var (
		structName     = bigCamel(table)
		thisName       = strings.ToLower(structName[:1])
		columns        = bytes.NewBuffer(nil)
		rows           = bytes.NewBuffer(nil)
		structContent  = bytes.NewBuffer(nil)
		structLabel    = bytes.NewBuffer(nil)
		structAssemble = bytes.NewBuffer(nil)
		fieldMap       = bytes.NewBuffer(nil)
	)

	if primaryField == "" {
		primaryField = "id"
	}

	for i, c := range cols {
		var (
			fieldName    = bigCamel(c.Field)
			propertyName = smallCamel(c.Field)
			tag          = fmt.Sprintf("`json:\"%s\"`", propertyName)
			comment      = ""
		)

		if c.Comment != "" {
			comment = " // " + c.Comment
			if c.DefaultVal != nil {
				comment += fmt.Sprintf(" (default: %s)", *c.DefaultVal)
			}
		} else if c.DefaultVal != nil {
			comment = fmt.Sprintf(" // default: %s", *c.DefaultVal)
		}

		if i > 0 {
			structLabel.WriteString("\n\t\t")
			structContent.WriteString("\n\t")
			structAssemble.WriteString("\n\t")
			fieldMap.WriteString("\n\t\t")
		}

		if columns.Len() > 0 {
			columns.WriteString(",")
			rows.WriteString(",")
		}

		if c.Field != primaryField {
			columns.WriteString(fmt.Sprintf(`"%s"`, c.Field))
			rows.WriteString(fmt.Sprintf(`info.%s`, fieldName))
		}

		structLabel.WriteString(fmt.Sprintf("\"%s\": \"%s\",", propertyName, c.Comment))
		fieldMap.WriteString(fmt.Sprintf("\"%s\": \"%s\",", propertyName, c.Field))
		structContent.WriteString(fmt.Sprintf("%-20s %-15s %-30s%s", fieldName, c.GoType, tag, comment))
		structAssemble.WriteString(fmt.Sprintf("%s.%s = ", thisName, fieldName))
		switch c.GoType {
		case "string":
			structAssemble.WriteString(fmt.Sprintf("br.String(\"%s\")", c.Field))
		case "[]byte":
			structAssemble.WriteString(fmt.Sprintf("br.Bytes(\"%s\")", c.Field))
		default:
			structAssemble.WriteString(fmt.Sprintf("br.To%s(\"%s\")", strings.ToUpper(c.GoType[:1])+c.GoType[1:], c.Field))
		}
	}

	model = bytes.ReplaceAll(ModelTemplate(), []byte("{pkg}"), []byte(pkg))
	model = bytes.ReplaceAll(model, []byte("{tableName}"), []byte(table))
	model = bytes.ReplaceAll(model, []byte("{primaryField}"), []byte(primaryField))
	model = bytes.ReplaceAll(model, []byte("{structName}"), []byte(structName))
	model = bytes.ReplaceAll(model, []byte("{this}"), []byte(thisName))
	model = bytes.ReplaceAll(model, []byte("{structContent}"), structContent.Bytes())
	model = bytes.ReplaceAll(model, []byte("{fieldMap}"), fieldMap.Bytes())
	model = bytes.ReplaceAll(model, []byte("{structLabel}"), structLabel.Bytes())
	model = bytes.ReplaceAll(model, []byte("{structAssemble}"), structAssemble.Bytes())

	crud = bytes.ReplaceAll(CrudTemplate(), []byte("{pkg}"), []byte(pkg))
	crud = bytes.ReplaceAll(crud, []byte("{primaryField}"), []byte(primaryField))
	crud = bytes.ReplaceAll(crud, []byte("{structName}"), []byte(structName))
	crud = bytes.ReplaceAll(crud, []byte("{this}"), []byte(thisName))
	crud = bytes.ReplaceAll(crud, []byte("{rows}"), rows.Bytes())
	crud = bytes.ReplaceAll(crud, []byte("{columns}"), columns.Bytes())
	return
}
