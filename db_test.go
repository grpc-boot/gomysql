package gomysql

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

var (
	db *Db
)

func init() {
	var err error
	db, err = NewDb(Options{
		DbName:   "users",
		Password: "12345678",
	})

	if err != nil {
		log.Fatalf("init db failed with error: %v\n", err)
	}
}

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

func TestDb_Find(t *testing.T) {
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
}
