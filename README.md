# gomysql
go语言实现的mysql帮助库

### 特性

> 1.基于泛型，无反射

> 2.类似Sql语法，语法更清晰明了

> 3.支持读写分离

> 4.支持节点failover和failback

> 5.支持全局sql输出到日志，自动记录错误日志

> 6.gmm代码生成器

<!-- TOC -->

- 使用说明
  - [实例化db](#实例化db) 
  - [Create-Table](#Create-Table) 
  - [Model](#Model)
  - [Select](#Select) 
  - [Insert](#Insert) 
  - [Update](#Update) 
  - [Aggregation](#Aggregation) 
  - [Delete](#Delete) 
  - [Transaction](#Transaction)
  - [Read-Write-Splitting](#Read-Write-Splitting) 
  - [Sql-Log](#Sql-Log) 
  - [测试](db_test.go)

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

#### Create-Table

```go
package main

import (
	"testing"

    "github.com/grpc-boot/gomysql"
	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_Exec(t *testing.T) {
	res, err := gomysql.Exec(db.Executor(), "CREATE TABLE IF NOT EXISTS `users` " +
		"(`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id'," +
		"`user_name` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '登录名'," +
		"`nickname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称'," +
		"`passwd` char(32) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL DEFAULT '' COMMENT '密码'," +
		"`email` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '邮箱'," +
		"`mobile` varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL DEFAULT '' COMMENT '手机号'," +
		"`is_on` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '账号状态(1已启用，0已禁用)'," +
		"`created_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间'," +
		"`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'," +
		"`last_login_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上次登录时间'," +
		"`remark` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '备注'," +
		"PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;")
	if err != nil {
		t.Fatalf("create table failed with error: %v\n", err)
	}

	count, err := res.RowsAffected()
	t.Logf("rows affected: %d error: %v\n", count, err)
}
```

#### Model

```go
package main

import (
  "fmt"
  "log"
  "time"

  "github.com/grpc-boot/gomysql"
  "github.com/grpc-boot/gomysql/helper"
)

var (
  DefaultUserModel = &UserModel{}
  db               *gomysql.Db
)

func init() {
  var err error
  db, err = gomysql.NewDb(gomysql.Options{
    Host:     "127.0.0.1",
    Port:     3306,
    DbName:   "users",
    UserName: "root",
    Password: "12345678",
  })

  if err != nil {
    log.Fatalf("init db failed with error: %v\n", err)
  }

  gomysql.SetLogger(func(query string, args ...any) {
    fmt.Printf("%s exec sql: %s args: %+v\n", time.Now().Format(time.DateTime), query, args)
  })

  gomysql.SetErrorLog(func(err error, query string, args ...any) {
    fmt.Printf("%s exec sql: %s args: %+v with error: %v\n", time.Now().Format(time.DateTime), query, args, err)
  })
}

func main() {
  current := time.Now()
  res, err := gomysql.Insert(
    db.Executor(),
    DefaultUserModel.TableName(),
    helper.Columns{"user_name", "nickname", "passwd", "is_on", "created_at", "updated_at", "last_login_at"},
    helper.Row{"uname", "nickName", "passwd", 1, current.Unix(), current.Format(time.DateTime), current.Format(time.DateTime)},
  )

  if err != nil {
    panic(err)
  }

  id, _ := res.LastInsertId()
  fmt.Printf("insert id: %d\n", id)

  user, err := gomysql.FindById(db.Pool(), id, DefaultUserModel)
  if err != nil {
    panic(err)
  }

  fmt.Printf("UserInfo: %+v\n", user)
}

type UserModel struct {
  Id          int64
  UserName    string
  NickName    string
  Passwd      string
  Email       string
  Mobile      string
  IsOn        uint8
  CreatedAt   int64
  UpdatedAt   time.Time
  LastLoginAt time.Time
  Remark      string
}

func (um *UserModel) PrimaryKey() string {
  return `id`
}

func (um *UserModel) TableName(args ...any) string {
  return `users`
}

func (um *UserModel) NewModel() gomysql.Model {
  return &UserModel{}
}

func (um *UserModel) Assemble(br gomysql.BytesRecord) {
  fmt.Printf("updated_at:%s\n", br.String("updated_at"))
  um.Id = br.ToInt64("id")
  um.UserName = br.String("user_name")
  um.NickName = br.String("nickname")
  um.Passwd = br.String("passwd")
  um.Email = br.String("email")
  um.Mobile = br.String("mobile")
  um.IsOn = br.ToUint8("is_on")
  um.CreatedAt = br.ToInt64("created_at")
  um.UpdatedAt, _ = time.Parse(time.DateTime, br.String("updated_at"))
  um.LastLoginAt, _ = time.Parse(time.DateTime, br.String("last_login_at"))
  um.Remark = br.String("remark")
}
```

#### Select

```go
package main

import (
	"testing"
	"time"

    "github.com/grpc-boot/gomysql"
	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_Find(t *testing.T) {
  // SELECT * FROM users WHERE id=1
  query := helper.AcquireQuery().
    From(`users`).
    Where(condition.Equal{"id", 2})

  defer query.Close()

  record, err := gomysql.FindOne(db.Executor(), query)
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }
  t.Logf("record: %+v\n", record)

  // SELECT * FROM users WHERE id IN(1, 2)
  query1 := helper.AcquireQuery().
    From(`users`).
    Where(condition.In[int]{"id", []int{1, 2}})
  defer query1.Close()

  records, err := gomysql.FindTimeout(time.Second*2, db.Executor(), query1)
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }
  t.Logf("records: %+v\n", records)

  // SELECT * FROM users WHERE user_name LIKE 'user%' AND created_at> timestamp
  query2 := helper.AcquireQuery().
    From(`users`).
          Where(condition.And{
            condition.BeginWith{"user_name", "user"},
            condition.Gte{"created_at", time.Now().Add(-7 * 24 * time.Hour).Unix()},
          })
  defer query2.Close()

  records, err = gomysql.FindTimeout(time.Second*2, db.Executor(), query2)
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

  records, err = gomysql.FindTimeout(time.Second*2, db.Executor(), query3)
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }
  t.Logf("records: %+v\n", records)
}

func TestFindModel(t *testing.T) {
  // SELECT * FROM users WHERE id=? args: [2]
  u, err := gomysql.FindModelByIdTimeout(time.Second, DefaultUserModel, db.Executor(), 2)
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }
  t.Logf("user: %+v\n", u)

  // SELECT * FROM users WHERE id=? args: [2]
  u, err = gomysql.FindModelByConditionTimeout(time.Second, DefaultUserModel, db.Executor(), condition.Equal{Field: "id", Value: 2})
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }
  t.Logf("user: %+v\n", u)

  // SELECT * FROM users WHERE id<? args: [100]
  us, err := gomysql.FindModelsByConditionTimeout(time.Second, DefaultUserModel, db.Executor(), condition.Lt{Field: "id", Value: 100})
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }
  t.Logf("users: %#v\n", us)

  q := helper.AcquireQuery().
    From(DefaultUserModel.TableName()).
    Where(condition.Gt{Field: "id", Value: 1}).
    Limit(100)
  defer q.Close()
  // SELECT * FROM users WHERE id>? LIMIT 0,100 args: [1]
  us, err = gomysql.FindModelsTimeout(time.Second, DefaultUserModel, db.Executor(), q)
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }
  t.Logf("users: %#v\n", us)

  // SELECT * FROM users WHERE id>? LIMIT 0,100 args: [1]
  u, err = gomysql.FindModelTimeout(time.Second, DefaultUserModel, db.Executor(), q)
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }
  t.Logf("user: %+v\n", u)
}
```

#### Insert

```go
package main

import (
	"strings"
	"testing"
	"time"

    "github.com/grpc-boot/gomysql"
	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func TestDb_Insert(t *testing.T) {
  res, err := gomysql.Insert(
    db.Executor(),
    `users`,
    helper.Columns{"user_name", "nickname", "passwd", "is_on", "created_at", "updated_at"},
    helper.Row{"user1", "nickname1", strings.Repeat("1", 32), 1, time.Now().Unix(), time.Now().Format(time.DateTime)},
  )

  if err != nil {
    t.Fatalf("insert data failed with error: %v\n", err)
  }

  id, _ := res.LastInsertId()
  t.Logf("insert data with id: %d\n", id)

  id, err = gomysql.InsertWithInsertedIdTimeout(
    time.Second,
    db.Executor(),
    `users`,
    helper.Columns{"user_name", "nickname", "passwd", "is_on", "created_at", "updated_at"},
    helper.Row{"user1", "nickname1", strings.Repeat("1", 32), 1, time.Now().Unix(), time.Now().Format(time.DateTime)},
  )
  if err != nil {
    t.Fatalf("insert data failed with error: %v\n", err)
  }

  t.Logf("insert data with id: %d\n", id)
}
```

#### Update

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

func TestDb_Update(t *testing.T) {
  res, err := gomysql.Update(
    db.Executor(),
    `users`,
    `last_login_at=?`,
    condition.Equal{Field: "id", Value: 2},
    time.Now().Format(time.DateTime),
  )

  if err != nil {
    t.Fatalf("update data failed with error: %v\n", err)
  }

  rows, _ := res.RowsAffected()
  t.Logf("update data rows affected: %d", rows)

  rows, err = gomysql.UpdateWithRowsAffectedContext(
    context.Background(),
    db.Executor(),
    `users`,
    `last_login_at=?`,
    condition.Equal{Field: "id", Value: 3},
    time.Now().Format(time.DateTime),
  )

  if err != nil {
    t.Fatalf("update data failed with error: %v\n", err)
  }

  t.Logf("update data rows affected: %d", rows)
}
```

#### Aggregation

```go
package main

import (
  "testing"
  "time"

  "github.com/grpc-boot/gomysql"
  "github.com/grpc-boot/gomysql/condition"
  "github.com/grpc-boot/gomysql/helper"
)

func Test_Agg(t *testing.T) {
	count, err := gomysql.CountByConditionTimeout(
		time.Second,
        DefaultUserModel,
		db.Executor(),
		condition.Gt{Field: "id", Value: 0},
		"id",
	)

	if err != nil {
		t.Fatalf("count failed with error: %v\n", err)
	}
	t.Logf("count: %d", count)

	sum, err := gomysql.SumByConditionTimeout(
		time.Second,
        DefaultUserModel,
		db.Executor(),
		condition.Gt{Field: "id", Value: 0},
		"id",
	)

	if err != nil {
		t.Fatalf("sum failed with error: %v\n", err)
	}
	t.Logf("sum: %d", int64(sum))

	maxId, err := gomysql.MaxByConditionTimeout(
		time.Second,
		defaultUser,
		db.Executor(),
		condition.Gt{Field: "id", Value: 0},
		"id",
	)

	if err != nil {
		t.Fatalf("max failed with error: %v\n", err)
	}
	t.Logf("maxId: %s", maxId)

	minId, err := gomysql.MinByConditionTimeout(
		time.Second,
		defaultUser,
		db.Executor(),
		condition.Gt{Field: "id", Value: 0},
		"id",
	)

	if err != nil {
		t.Fatalf("min failed with error: %v\n", err)
	}
	t.Logf("minId: %s", minId)

	avg, err := gomysql.AvgByConditionTimeout(
		time.Second,
		defaultUser,
		db.Executor(),
		condition.Gt{Field: "id", Value: 0},
		"id",
	)

	if err != nil {
		t.Fatalf("avg failed with error: %v\n", err)
	}
	t.Logf("avg: %s", avg)
}
```

#### Delete

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

func TestDb_Delete(t *testing.T) {
  res, err := gomysql.Delete(db.Executor(), `users`, condition.Equal{Field: "id", Value: 1})
  if err != nil {
    t.Fatalf("delete data failed with error: %v\n", err)
  }

  rows, _ := res.RowsAffected()
  t.Logf("delete data rows affected: %d", rows)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  rows, err = gomysql.DeleteWithRowsAffectedContext(
    ctx,
    db.Executor(),
    `users`,
    condition.Equal{Field: "id", Value: 1},
  )

  if err != nil {
    t.Fatalf("delete data failed with error: %v\n", err)
  }
  t.Logf("delete data rows affected: %d", rows)
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

  defer tx.Rollback()
  
  query := helper.AcquireQuery().
    From(`users`).
    Where(condition.Equal{"id", 1})
  defer query.Close()
  records, err := gomysql.Find(tx, query)
  if err != nil {
    t.Fatalf("query failed with error: %v", err)
  }

  if len(records) != 1 {
    t.Fatal("row not exists")
  }

  res, err := gomysql.Update(tx, `users`, "updated_at=?", condition.Equal{"updated_at", records[0].String("updated_at")}, time.Now().Format(time.DateTime))
  if err != nil {
    t.Fatalf("update failed with error: %v", err)
  }

  tx.Commit()
  count, _ := res.RowsAffected()
  t.Logf("updated count: %d", count)
}
```

#### Read-Write-Splitting

> 支持failover和failback

```go
package main

import (
  "testing"
  "time"
  
  "github.com/grpc-boot/gomysql"
  "github.com/grpc-boot/gomysql/condition"
  "github.com/grpc-boot/gomysql/helper"
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

  pool, err := gomysql.NewPool(opt, time.Second*10)
  if err != nil {
    t.Fatalf("want nil, got %v", err)
  }

  var (
    query = helper.AcquireQuery().
      From(`users`).
      Where(condition.Equal{"id", 1})
    start       = time.Now()
    maxInterval = time.Minute
  )

  defer query.Close()

  record, err := gomysql.FindOne(pool.RandExecutor(gomysql.TypeMaster), query)
  if err != nil {
    t.Logf("find one error: %v", err)
  } else {
    t.Logf("query records: %+v", record)
  }

  ticker := time.NewTicker(time.Second * 5)
  for range ticker.C {
    if time.Since(start) > maxInterval {
      ticker.Stop()
      break
    }

    record, err = gomysql.FindOne(pool.RandExecutor(gomysql.TypeSlave), query)
    if err != nil {
      t.Logf("find one error: %v", err)
    } else {
      t.Logf("query records: %+v", record)
    }
  }
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

  // 记录错误日志
  gomysql.SetErrorLog(func(err error, query string, args ...any) {
    fmt.Printf("%s exec sql: %s args: %+v with error: %v\n", time.Now().Format(time.DateTime), query, args, err)
  })
}

```

#### gmm代码生成器

```shell
# 安装
$ go install github.com/grpc-boot/gomysql/example/gmm@v1.1.6

# 使用帮助
$ ./gmm
Usage of ./gmm:
  -P int
    	Mysql port (default 3306)
  -c string
    	Charset (default "utf8")
  -d string
    	Mysql db name
  -h string
    	Mysql host (default "localhost")
  -o string
    	Gmm model output directory (default "entity")
  -p string
    	Mysql password
  -t string
    	Gmm table, All tables are selected by default
  -u string
    	Mysql user (default "root")
    	
# 生成代码
$ ./gmm -d=users
create model for table:users success
create model for table:users_log success
```

