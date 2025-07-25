package gomysql

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/filter"
	"github.com/grpc-boot/gomysql/helper"
)

var (
	db          *Db
	defaultUser = &User{}
)

type User struct {
	Id       int64  `json:"id"`
	UserName string `json:"userName"`
}

func (u *User) PrimaryKey() string {
	return `id`
}

func (u *User) TableName(args ...any) string {
	return `users`
}

func (u *User) NewModel() Model {
	return &User{}
}

func (u *User) Assemble(br BytesRecord) {
	u.Id = br.ToInt64("id")
	u.UserName = br.String("user_name")
}

func init() {
	var err error
	db, err = NewDb(Options{
		Host:     "127.0.0.1",
		Port:     3306,
		DbName:   "users",
		UserName: "root",
		Password: "",
	})

	if err != nil {
		log.Fatalf("init db failed with error: %v\n", err)
	}

	SetLogger(func(query string, args ...any) {
		fmt.Printf("%s exec sql: %s args: %+v\n", time.Now().Format(time.DateTime), query, args)
	})

	SetErrorLog(func(err error, query string, args ...any) {
		fmt.Printf("%s exec sql: %s args: %+v with error: %v\n", time.Now().Format(time.DateTime), query, args, err)
	})
}

func TestDb_Exec(t *testing.T) {
	res, err := Exec(db.Executor(), "CREATE TABLE IF NOT EXISTS `users` "+
		"(`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',"+
		"`user_name` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '登录名',"+
		"`nickname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',"+
		"`passwd` char(32) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL DEFAULT '' COMMENT '密码',"+
		"`email` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',"+
		"`mobile` varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL DEFAULT '' COMMENT '手机号',"+
		"`is_on` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '账号状态(1已启用，0已禁用)',"+
		"`created_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',"+
		"`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',"+
		"`last_login_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上次登录时间',"+
		"`remark` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '备注',"+
		"PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;")
	if err != nil {
		t.Fatalf("create table failed with error: %v\n", err)
	}

	count, err := res.RowsAffected()
	t.Logf("rows affected: %d error: %v\n", count, err)
}

func TestPrepare(t *testing.T) {
	_, err := Prepare(db.Executor(), "SELECT * FROM users")
	if err != nil {
		t.Fatalf("prepare query failed with error: %v\n", err)
	}
}

func TestScrollFilterWithIdContext(t *testing.T) {
	var caseList = []struct {
		name string
		f    filter.Scroll
	}{
		{
			name: "id降序",
			f: filter.Scroll{
				PageSize: 30,
			},
		},
		{
			name: "id升序，大于等于34",
			f: filter.Scroll{
				Filters: map[string]filter.Filter{},
				Sorts: filter.SortItem{
					Field: "id",
				},
				Cursor:   "34",
				PageSize: 30,
			},
		},
	}

	for _, cs := range caseList {
		t.Run(cs.name, func(t *testing.T) {
			hasMore, last, rows, err := ScrollFilterWithIdTimeout(time.Second, &cs.f, math.MaxInt64, defaultUser, db.Executor())
			t.Logf("hasMore: %v last: %v rows: %v error: %v", hasMore, last, rows, err)
			if err != nil {
				t.Fatalf("query failed with error: %v\n", err)
			}
		})
	}
}

func TestDb_Insert(t *testing.T) {
	res, err := Insert(
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

	id, err = InsertWithInsertedIdTimeout(
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

func TestDb_Update(t *testing.T) {
	res, err := Update(
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

	rows, err = UpdateWithRowsAffectedContext(
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

func Test_Agg(t *testing.T) {
	count, err := CountByConditionTimeout(
		time.Second,
		defaultUser,
		db.Executor(),
		condition.Gt{Field: "id", Value: 0},
		"id",
	)

	if err != nil {
		t.Fatalf("count failed with error: %v\n", err)
	}
	t.Logf("count: %d", count)

	sum, err := SumByConditionTimeout(
		time.Second,
		defaultUser,
		db.Executor(),
		condition.Gt{Field: "id", Value: 0},
		"id",
	)

	if err != nil {
		t.Fatalf("sum failed with error: %v\n", err)
	}
	t.Logf("sum: %d", int64(sum))

	maxId, err := MaxByConditionTimeout(
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

	minId, err := MinByConditionTimeout(
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

	avg, err := AvgByConditionTimeout(
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

func TestDb_Delete(t *testing.T) {
	res, err := Delete(db.Executor(), `users`, condition.Equal{Field: "id", Value: 1})
	if err != nil {
		t.Fatalf("delete data failed with error: %v\n", err)
	}

	rows, _ := res.RowsAffected()
	t.Logf("delete data rows affected: %d", rows)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rows, err = DeleteWithRowsAffectedContext(
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

func TestDb_Find(t *testing.T) {
	// SELECT * FROM users WHERE id=1
	query := helper.AcquireQuery().
		From(`users`).
		Where(condition.Equal{"id", 2})

	defer query.Close()

	record, err := FindOne(db.Executor(), query)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("record: %+v\n", record)

	// SELECT * FROM users WHERE id IN(1, 2)
	query1 := helper.AcquireQuery().
		From(`users`).
		Where(condition.In[int]{"id", []int{1, 2}})
	defer query1.Close()

	records, err := FindTimeout(time.Second*2, db.Executor(), query1)
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

	records, err = FindTimeout(time.Second*2, db.Executor(), query2)
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

	records, err = FindTimeout(time.Second*2, db.Executor(), query3)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("records: %+v\n", records)
}

func TestFindModel(t *testing.T) {
	// SELECT * FROM users WHERE id=? args: [2]
	u, err := FindModelByIdTimeout(time.Second, defaultUser, db.Executor(), 2)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("user: %+v\n", u)

	// SELECT * FROM users WHERE id=? args: [2]
	u, err = FindModelByConditionTimeout(time.Second, defaultUser, db.Executor(), condition.Equal{Field: "id", Value: 2})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("user: %+v\n", u)

	// SELECT * FROM users WHERE id<? args: [100]
	us, err := FindModelsByConditionTimeout(time.Second, defaultUser, db.Executor(), condition.Lt{Field: "id", Value: 100})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("users: %#v\n", us)

	q := helper.AcquireQuery().
		From(defaultUser.TableName()).
		Where(condition.Gt{Field: "id", Value: 1}).
		Limit(100)
	defer q.Close()
	// SELECT * FROM users WHERE id>? LIMIT 0,100 args: [1]
	us, err = FindModelsTimeout(time.Second, defaultUser, db.Executor(), q)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("users: %#v\n", us)

	// SELECT * FROM users WHERE id>? LIMIT 0,100 args: [1]
	u, err = FindModelTimeout(time.Second, defaultUser, db.Executor(), q)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("user: %+v\n", u)
}

func TestPool_Random(t *testing.T) {
	opt := PoolOptions{
		Masters: []Options{
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "",
			},
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "",
			},
		},
		Slaves: []Options{
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "",
			},
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "",
			},
			{
				Host:     "127.0.0.1",
				Port:     3306,
				DbName:   "users",
				UserName: "root",
				Password: "",
			},
		},
	}

	pool, err := NewPool(opt, time.Second*10)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	var (
		query = helper.AcquireQuery().
			From(`users`).
			Where(condition.Equal{"id", 1})
		start       = time.Now()
		maxInterval = time.Second * 20
	)

	defer query.Close()

	record, err := FindOne(pool.RandExecutor(TypeMaster), query)
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

		record, err = FindOne(pool.RandExecutor(TypeSlave), query)
		if err != nil {
			t.Logf("find one error: %v", err)
		} else {
			t.Logf("query records: %+v", record)
		}
	}
}

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
		Where(condition.Gte{"id", 1})
	defer query.Close()
	records, err := Find(tx, query)
	if err != nil {
		t.Fatalf("query failed with error: %v", err)
	}

	if len(records) != 1 {
		t.Fatal("row not exists")
	}

	res, err := Update(tx, `users`, "updated_at=?", condition.Equal{"updated_at", records[0].String("updated_at")}, time.Now().Format(time.DateTime))
	if err != nil {
		t.Fatalf("update failed with error: %v", err)
	}

	tx.Commit()
	count, _ := res.RowsAffected()
	t.Logf("updated count: %d", count)
}
