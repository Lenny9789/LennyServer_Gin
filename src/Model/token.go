package Model

import (
	"github.com/go-xorm/xorm"
	"log"
)

type Token struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	UserId int64 `json:"user_id"`
	TokenString string `json:"token_string"`
}

func init() {
	var err error
	//创建ORM引擎，链接数据库
	x, err = xorm.NewEngine("sqlite3", "./user.db")
	if err != nil {
		log.Fatalf("Fail to create engine: %v\n", err)
	}
	//同步结构体与数据表
	if err = x.Sync(new(Token)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}
}
//插入数据库
func (t *Token) SaveToDatabase() error {
	//对不存在的记录进行插入
	_, err := x.Insert(t)
	if err != nil {
		return err
	}
	return nil
}
//更新数据
func (t *Token) Update() error {
	//对已有记录进行更新
	_, err := x.Update(t)
	if err != nil {
		return err
	}
	return nil
}
//从数据库删除
func (t *Token) DeleteFromDatabase() error {
	_, err := x.Delete(t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) IsNull() bool {
	if t.Id == 0 || t.Username == "" || t.UserId == 0 || t.TokenString == "" {
		return true
	}else {
		return false
	}
}

//获取用户登录的Token信息
func Get_A_UserToken(userId int64) (*Token, error) {
	t := &Token{UserId:userId}
	err := x.Find(t)
	return t, err
}

func Get_A_UserToken_Info(t *Token) *Token {
	t2 := &Token{}
	x.Alias("t").Where("t.Username = ?", t.Username).Get(t2)
	if t2.IsNull()  {
		return nil
	}else {
		return t2
	}
}