package {pkg}

import (
    "context"
    "time"

    "github.com/grpc-boot/gomysql"
    "github.com/grpc-boot/gomysql/condition"
    "github.com/grpc-boot/gomysql/filter"
    "github.com/grpc-boot/gomysql/helper"
)

var (
    Default{structName} = &{structName}{}
)

// SearchByScroll
/** filter example:
* {
        "filters": {  // 多个字段间筛选时AND关系
            "field1": {
                "field": "field1",
                "opt": ">", //=、≠、>、>=、<、<=、Contains、Not Contains、Start With、End With、Is Not Empty、Is Empty、Is Any Of
                "value": "100", //当opt为Is Any Of时，多个value值通过","分隔
            },
            "field2": {
                "field": "field2",
                "opt": "=",
                "value": "2",
            }
        },
        "sorts": {     // cursor分页仅支持单字段排序
            "field": "field3",
            "kind": "desc" //asc、desc
        },
        "cursor": "0",
        "pageSize": 100,
    }
*/
func ({this} *{structName}) SearchByScroll(timeout time.Duration, filter *filter.Scroll, defaultId int64, db gomysql.Executor, args ...any) (hasMore bool, last *{structName}, ms []*{structName}, err error) {
    return gomysql.ScrollFilterWithIdTimeout(timeout, filter, defaultId, Default{structName}, db, args...)
}

// SearchByPage
/**  filter example:
*   {
        "filters": {
            "field1": {
                "field": "field1",
                "opt": ">", //=、≠、>、>=、<、<=、Contains、Not Contains、Start With、End With、Is Not Empty、Is Empty、Is Any Of
                "value": "100", //当opt为Is Any Of时，多个value值通过","分隔
            },
            "field2": {
                "field": "field2",
                "opt": "=",
                "value": "2",
            }
        },
        "sorts": [
            {
                "field": "field3",
                "kind": "desc" //asc、desc
            },
            {
                "field": "field4",
                "kind": "asc" //asc、desc
            },
        ],
        "page": 1,
        "pageSize": 100,
    }
*/
func ({this} *{structName}) SearchByPage(timeout time.Duration, filter *filter.Page, db gomysql.Executor, args ...any) (totalCount, pageCount int64, ms []*{structName}, err error) {
    return gomysql.PageFilterWithTotalTimeout(timeout, filter, Default{structName}, db, args...)
}

func ({this} *{structName}) Index(timeout time.Duration,db gomysql.Executor, cond condition.Condition, tableArgs ...any) (ms []*{structName}, err error){
    return gomysql.FindModelsByConditionTimeout(timeout, Default{structName}, db, cond, tableArgs...)
}

func ({this} *{structName}) Info(timeout time.Duration,db gomysql.Executor, cond condition.Condition, tableArgs ...any) (info *{structName}, err error){
    return gomysql.FindModelByConditionTimeout(timeout, Default{structName}, db, cond, tableArgs...)
}

func ({this} *{structName}) Create(db gomysql.Executor, info *{structName}) (id int64, err error) {
    return gomysql.InsertWithInsertedIdContext(
        context.Background(),
        db,
        Default{structName}.TableName(),
        helper.Columns{{columns}},
        helper.Row{{rows}},
    )
}

func ({this} *{structName}) InfoById(db gomysql.Executor, id int64) (info *{structName}, err error) {
    return gomysql.FindById(db, id, Default{structName})
}

func ({this} *{structName}) Update(db gomysql.Executor, id uint64, setters string, setterArgs ...any) (rows int64, err error) {
    return gomysql.UpdateWithRowsAffectedContext(
        context.Background(),
        db,
        Default{structName}.TableName(),
        setters,
        condition.Equal{Field: "{primaryField}", Value: id},
        setterArgs...,
    )
}

func ({this} *{structName}) Delete(db gomysql.Executor, id uint64) (rows int64, err error) {
    return gomysql.DeleteWithRowsAffectedContext(
        context.Background(),
        db,
        Default{structName}.TableName(),
        condition.Equal{Field: "{primaryField}", Value: id},
    )
}