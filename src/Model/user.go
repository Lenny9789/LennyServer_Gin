package Model

import (
	"github.com/go-xorm/xorm"
	"log"
	"github.com/kataras/iris/core/errors"
)

type User struct {
	Id int64 `json:"id"`
	Username string `xorm:"unique", json:"username"`
	Password string `json:"password"`
	Phonenumber string `json:"phonenumber"`
	Realname string `json:"realname"`
	Age int8 `json:"age"`
	Sex int8 `json:"sex"`
}


func init() {
	//创建ORM引擎与数据库
	var err error
	x, err = xorm.NewEngine("sqlite3", "./user.db")
	if err != nil {
		log.Fatalf("Fail to create engine: %v\n", err)
	}
	//同步结构体与数据表
	if err = x.Sync(new(User)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}
}
//插入数据库
func (u *User) SaveToDatabase() error {
	//对不存在的记录进行插入
	_, err := x.Insert(u)
	if err != nil {
		return err
	}
	return nil
}

//更新
func (u *User) Update() error {
	//对已有记录进行更新
	_, err := x.Update(u)
	if err != nil {
		return err
	}
	return nil
}
//从数据库删除
func (u *User) DeleteFromDatabase() error {
	_, err := x.Delete(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) IsNull() bool {
	if u.Username == "" || u.Password == "" || u.Id == 0 {
		return true
	}else {
		return false
	}
}

//获取User信息
func Get_A_User(id int64) (*User, error) {
	u := &User{}
	//直接操作ID的简便方法
	has, err := x.Id(id).Get(u)
	if err != nil {
		return nil, err
	}else if !has {
		return nil, errors.New("User does not exsit.")
	}
	return u, nil
}
//获取User信息
func Get_A_User_Info(u *User) *User {
	u2 := &User{}
	x.Alias("u").Where("u.Username = ?", u.Username).Get(u2)
	if u2.IsNull() {
		return nil
	}else {
		return u2
	}
}
//按照ID正排序返回所有User
func Get_Users_AscId() (as []User, err error) {
	//使用Find方法批量获取记录
	err = x.Find(&as)
	return as, err
}

