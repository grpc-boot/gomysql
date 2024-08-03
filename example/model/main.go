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
}

func main() {
	current := time.Now()
	res, err := db.Insert(
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
