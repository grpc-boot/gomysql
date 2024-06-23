package gomysql

import (
	"fmt"
	"time"
)

func DefaultMysqlOption() Options {
	return Options{
		Host:                  "127.0.0.1",
		Port:                  3306,
		UserName:              "root",
		CharSet:               "utf8",
		MaxOpenConns:          8,
		MaxIdleConns:          2,
		ConnMaxLifetimeSecond: 600,
		ConnMaxIdleTimeSecond: 60,
	}
}

type Options struct {
	DbName                string `json:"dbName" yaml:"dbName"`
	Host                  string `json:"host" yaml:"host"`
	Port                  uint32 `json:"port" yaml:"port"`
	UserName              string `json:"userName" yaml:"userName"`
	Password              string `json:"password" yaml:"password"`
	CharSet               string `json:"charSet" yaml:"charSet"`
	MaxIdleConns          int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns          int    `json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxIdleTimeSecond int64  `json:"connMaxIdleTimeSecond" yaml:"connMaxIdleTimeSecond"`
	ConnMaxLifetimeSecond int64  `json:"connMaxLifetimeSecond" yaml:"connMaxLifetimeSecond"`
}

func (o *Options) format() *Options {
	defaultOpt := DefaultMysqlOption()
	if o.Host == "" {
		o.Host = defaultOpt.Host
	}

	if o.Port == 0 {
		o.Port = defaultOpt.Port
	}

	if o.UserName == "" {
		o.UserName = defaultOpt.UserName
	}

	if o.CharSet == "" {
		o.CharSet = defaultOpt.CharSet
	}

	if o.MaxOpenConns == 0 {
		o.MaxOpenConns = defaultOpt.MaxOpenConns
	}

	if o.ConnMaxLifetimeSecond == 0 {
		o.ConnMaxLifetimeSecond = defaultOpt.ConnMaxLifetimeSecond
	}

	if o.MaxIdleConns == 0 {
		o.MaxIdleConns = defaultOpt.MaxIdleConns
	}

	if o.ConnMaxIdleTimeSecond == 0 {
		o.ConnMaxIdleTimeSecond = defaultOpt.ConnMaxIdleTimeSecond
	}

	return o
}

func (o *Options) ConnMaxIdleTime() time.Duration {
	return time.Duration(o.ConnMaxIdleTimeSecond) * time.Second
}

func (o *Options) ConnMaxLifetime() time.Duration {
	return time.Duration(o.ConnMaxLifetimeSecond) * time.Second
}

func (o *Options) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		o.UserName,
		o.Password,
		o.Host,
		o.Port,
		o.DbName,
		o.CharSet,
	)
}
