package {pkg}

import (
    "github.com/grpc-boot/gomysql"
)

type {structName} struct {
    {structContent}
}

func ({this} *{structName}) NewModel() gomysql.Model {
    return &{structName}{}
}

func ({this} *{structName}) TableName(args ...any) string {
    return "{tableName}"
}

func ({this} *{structName}) LabelMap() map[string]string {
    return map[string]string{
        {structLabel}
    }
}

func ({this} *{structName}) PrimaryKey() string {
    return "{primaryField}"
}

func ({this} *{structName}) Assemble(br gomysql.BytesRecord) {
    if len(br) == 0 {
        return
    }

    {structAssemble}
}
