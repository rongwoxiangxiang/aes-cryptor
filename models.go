package main

import (
	"encoding/base64"
	"errors"
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

func (usr *User) GetByName(name string) (*User, error) {
	user := new(User)
	err := engin.NewOrm().Table(user).Where("name", name).Select()
	return user, err
}

func (usr *User) Login(name, pass string) error {
	user := new(User)
	err := engin.NewOrm().Table(user).Where("name", name).Select()
	if err != nil {
		return err
	} else if IsEmpty(user.Pass) {
		return errors.New("登录失败[0]")
	}
	signCalc := base64.StdEncoding.EncodeToString([]byte(pass + user.Salt))
	if Md5(signCalc) != user.Pass {
		return errors.New("密码错误")
	}
	return nil
}

type Log struct {
	Id        int64     `gorose:"id"`
	Operator  string    `gorose:"operator"`
	Content   string    `gorose:"content"`
	Operation string    `gorose:"operation"`
	CreatedOn time.Time `gorose:"created_on"`
}

func (lg *Log) Insert(log Log) (int64, error) {
	return engin.NewOrm().Table(&log).Data(log).Insert()
}
