package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	"time"
)

var engin *gorose.Engin

func init() {
	dsn := "root:root@tcp(localhost:3306)/crypt?charset=utf8&parseTime=true"
	localEngin, err := gorose.Open(&gorose.Config{Driver: "mysql", Dsn: dsn})
	if err != nil {
		fmt.Printf("Connect server err: %s", err.Error())
	}
	engin = localEngin
}

type User struct {
	Id        int64     `gorose:"id"`
	Name      string    `gorose:"name"`
	Pass      string    `gorose:"pass"`
	Salt      string    `gorose:"salt"`
	CreatedOn time.Time `gorose:"created_on"`
}

// 设置表名, 如果没有设置, 默认使用struct的名字
func (usr *User) TableName() string {
	return "user"
}

func (usr *User) GetByName(name string) (*User, error) {
	user := new(User)
	err := engin.NewOrm().Table(user).Where("name", name).Select()
	return user, err
}

type Log struct {
	Id        int64     `gorose:"id"`
	Operator  string    `gorose:"operator"`
	Content   string    `gorose:"content"`
	Operation string    `gorose:"operation"`
	CreatedOn time.Time `gorose:"created_on"`
}

func (lg *Log) TableName() string {
	return "log"
}

func (lg *Log) Insert(log Log) (int64, error) {
	return engin.NewOrm().Table(&log).Data(log).Insert()
}
