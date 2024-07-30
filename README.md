# gomysql
go语言实现的mysql帮助库

<!-- TOC -->

- 使用说明
  - [实例化db](#实例化db) 
  - [create-table](#create-table) 
  - [Select](#Select) 
  - [Insert](#Insert) 
  - [Update](#Update) 
  - [Delete](#Delete) 
  - [Transaction](#Transaction) 
  - [Model](#Model)
  - [Read-Write-Splitting](#Read-Write-Splitting) 
  - [Sql-Log](#Sql-Log) 

#### 实例化db

```go
package main

import (
	"log"
	
	"github.com/grpc-boot/gomysql"
)

func main() {
	db, err := gomysql.NewDb(gomysql.Options{
		Host:     "127.0.0.1",
		Port:     3306,
		DbName:   "users",
		UserName: "root",
		Password: "12345678",
	})

	if err != nil {
		log.Fatalf("init db failed with error: %v\n", err)
	}
}

```

#### create-table

```go
package main

import (
	"testing"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_Exec(t *testing.T) {
	res, err := db.Exec("CREATE TABLE IF NOT EXISTS `users` " +
		"(`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id'," +
		"`user_name` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '登录名'," +
		"`nickname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称'," +
		"`passwd` char(32) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL DEFAULT '' COMMENT '密码'," +
		"`email` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '邮箱'," +
		"`mobile` varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL DEFAULT '' COMMENT '手机号'," +
		"`is_on` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '账号状态(1已启用，0已禁用)'," +
		"`created_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间'," +
		"`updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间'," +
		"`last_login_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '上次登录时间'," +
		"`remark` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '备注'," +
		"PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;")
	if err != nil {
		t.Fatalf("create table failed with error: %v\n", err)
	}

	count, err := res.RowsAffected()
	t.Logf("rows affected: %d error: %v\n", count, err)
}
```

#### Select

```go
package main

import (
	"testing"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_Find(t *testing.T) {
	// SELECT * FROM users WHERE id=1
	query := helper.AcquireQuery().
		From(`users`).
		Where(condition.Equal{"id", 1})

	defer query.Close()

	record, err := db.FindOne(query)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("record: %+v\n", record)

	// SELECT * FROM users WHERE id IN(1, 2)
	query1 := helper.AcquireQuery().
		From(`users`).
		Where(condition.In[int]{"id", []int{1, 2}})
	defer query1.Close()

	records, err := db.FindTimeout(time.Second*2, query1)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("records: %+v\n", records)

	// SELECT * FROM users WHERE user_name LIKE 'user%' AND created_at>= timestamp
	query2 := helper.AcquireQuery().
		From(`users`).
		Where(condition.And{
			condition.BeginWith{"user_name", "user"},
			condition.Gte{"created_at", time.Now().Add(-7 * 24 * time.Hour).Unix()},
		})
	defer query2.Close()

	records, err = db.FindTimeout(time.Second*2, query2)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("records: %+v\n", records)

	// SELECT * FROM users WHERE (user_name LIKE 'animal' AND created_at BETWEEN timestamp1 AND timestamp2) OR user_name LIKE 'user%'
	query3 := helper.AcquireQuery().
		From(`users`).
		Where(condition.Or{
			condition.And{
				condition.BeginWith{"user_name", "animal"},
				condition.Between{"created_at", time.Now().Add(-30 * 7 * 24 * time.Hour).Unix(), time.Now().Unix()},
			},
			condition.BeginWith{"user_name", "user"},
		})
	defer query3.Close()

	records, err = db.FindTimeout(time.Second*2, query3)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("records: %+v\n", records)
}
```

#### Insert

```go
package main

import (
	"strings"
	"testing"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_Insert(t *testing.T) {
	res, err := db.Insert(
		`users`,
		helper.Columns{"user_name", "nickname", "passwd", "is_on", "created_at", "updated_at"},
		helper.Row{"user1", "nickname1", strings.Repeat("1", 32), 1, time.Now().Unix(), time.Now().Format(time.DateTime)},
	)
	if err != nil {
		t.Fatalf("insert data failed with error: %v\n", err)
	}

	id, _ := res.LastInsertId()
	t.Logf("insert data with id: %d\n", id)
}
```

#### Update

```go
package main

import (
	"testing"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_Update(t *testing.T) {
	res, err := db.UpdateTimeout(
		time.Second,
		`users`,
		`last_login_at=?`,
		condition.Equal{Field: "id", Value: 1},
		time.Now().Format(time.DateTime),
	)

	if err != nil {
		t.Fatalf("update data failed with error: %v\n", err)
	}

	count, _ := res.RowsAffected()
	t.Logf("update data rows affected: %d", count)
}
```

#### Delete

```go
package main

import (
	"context"
	"testing"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_Delete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := db.DeleteContext(
		ctx,
		`users`,
		condition.Equal{Field: "id", Value: 1},
	)

	if err != nil {
		t.Fatalf("delete data failed with error: %v\n", err)
	}

	count, _ := res.RowsAffected()
	t.Logf("delete data rows affected: %d", count)
}
```

#### Transaction

```go
package main

import (
	"context"
	"testing"
	"time"
	
	"github.com/grpc-boot/gomysql"
	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_BeginTx(t *testing.T) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
  defer cancel()
  tx, err := db.BeginTx(ctx, nil)
  if err != nil {
    t.Fatalf("begin failed with error: %v", err)
  }

  query := helper.AcquireQuery().
    From(`users`).
    Where(condition.Equal{"id", 1})
  defer query.Close()
  records, err := gomysql.Find(tx, query)
  if err != nil {
    tx.Rollback()
    t.Fatalf("query failed with error: %v", err)
  }

  if len(records) != 1 {
    tx.Rollback()
    t.Fatal("row not exists")
  }

  res, err := gomysql.Update(tx, `users`, "updated_at=?", condition.Equal{"updated_at", records[0].String("updated_at")}, time.Now().Format(time.DateTime))
  if err != nil {
    tx.Rollback()
    t.Fatalf("update failed with error: %v", err)
  }

  tx.Commit()
  count, _ := res.RowsAffected()
  t.Logf("updated count: %d", count)
}
```

#### Model

```go
package main

import (
  "testing"
  
  "github.com/grpc-boot/gomysql"
  "github.com/grpc-boot/gomysql/condition"
  "github.com/grpc-boot/gomysql/helper"
)

var (
	db *gomysql.Db
	DefaultUserModel = &UserModel{}
)

type UserModel struct {
	Id       int64
	UserName string
	Passwd   string
}

func (um *UserModel) NewModel() gomysql.Model {
	return &UserModel{}
}

func (um *UserModel) Assemble(br gomysql.BytesRecord) {
	um.Id = br.ToInt64("id")
	um.UserName = br.String("user_name")
	um.Passwd = br.String("passwd")
}

func TestBytesRecords2Model(t *testing.T) {
  brs := []gomysql.BytesRecord{
    {
      "id":        []byte(`100008834`),
      "user_name": []byte(`慌了神`),
      "passwd":    []byte(`sdfw9df239sadfj239fasdfadf`),
    },
    {
      "id":        []byte(`1000123834`),
      "user_name": []byte(`慌了神123dfasdf`),
      "passwd":    []byte(`adfj239fasdfadfasdfasfd`),
    },
  }

  models := gomysql.BytesRecords2Models(brs, DefaultUserModel)
  if len(models) != len(brs) {
    t.Fatalf("want %d, got %d", len(brs), len(models))
  }

  t.Logf("model 0: %+v", models[0])
}

func TestFindModel(t *testing.T) {
  query := helper.AcquireQuery().
    From(`users`).
    Where(condition.Equal{"id", 1})
  defer query.Close()

  users, err := gomysql.FindModels(DefaultUserModel, db.Pool(), query)
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }

  t.Logf("users: %+v", users)
}

```

#### Read-Write-Splitting

```go
package main

import (
  "testing"

  "github.com/grpc-boot/gomysql"
  "github.com/grpc-boot/gomysql/condition"
)

func TestPool_Random(t *testing.T) {
	opt := gomysql.PoolOptions{
		Masters: []gomysql.Options{
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "12345678",
			},
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "12345678",
			},
		},
		Slaves: []gomysql.Options{
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "12345678",
			},
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "12345678",
			},
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "12345678",
			},
		},
	}

	pool, err := gomysql.NewPool(opt)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	query := helper.AcquireQuery().
		From(`users`).
		Where(condition.Equal{"id", 1})
	defer query.Close()

	record, err := pool.Random(TypeMaster).FindOne(query)
	if err != nil {
		t.Fatalf("find one error: %v", err)
	}
	t.Logf("query records: %+v", record)

	record, err = pool.Random(TypeSlave).FindOne(query)
	if err != nil {
		t.Fatalf("find one error: %v", err)
	}
	t.Logf("query records: %+v", record)
}
```

#### Sql-Log

```go
package main

import (
	"fmt"
	"time"

	"github.com/grpc-boot/gomysql"
)

func init() {
  // 输出sql和参数到标准输入，修改func定制自己的日志，方便分析sql
  gomysql.SetLogger(func(query string, args ...any) {
    fmt.Printf("%s exec sql: %s args: %+v\n", time.Now().Format(time.DateTime), query, args)
  })
}

```

