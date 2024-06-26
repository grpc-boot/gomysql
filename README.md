# gomysql
go语言实现的mysql帮助库

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

#### Select

```go
// SELECT * FROM users WHERE id=1
query := helper.AcquireQuery().
    From(`users`).
    Where(condition.Equal{"id", 2})

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

// SELECT * FROM users WHERE user_name LIKE 'user%' AND created_at> timestamp
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
```
