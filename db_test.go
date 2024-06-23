package gomysql

import (
	"context"
	"fmt"
	"github.com/grpc-boot/gomysql/condition"
	"log"
	"strings"
	"testing"
	"time"

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
		fmt.Sprintf(`last_login_at="%s"`, time.Now().Format(time.DateTime)),
		condition.Equal{Field: "id", Value: 1},
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
