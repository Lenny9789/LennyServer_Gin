package Model

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/go-xorm/xorm"
	"log"
	"github.com/kataras/iris/core/errors"
)

type Account struct {
	Id int64
	Name string `xorm:"unique"`
	Balance float64
	Version int `xorm:"version"`
}

//var x *xorm.Engine

func init() {
	//创建ORM引擎与数据库
	var err error
	x, err = xorm.NewEngine("sqlite3", "./bank.db")
	if err != nil {
		log.Fatal("Fail to create engine: %v\n", err)
	}
	//同步结构体数据表
	if err = x.Sync(new(Account)); err != nil {
		log.Fatal("Fail to sync database: %v\n", err)
	}
}

func newAccount(name string, balance float64) error {
	//对不存在的记录进行插入
	_, err := x.Insert(&Account{Name:name, Balance:balance})
	return err
}
//获取账户信息
func getAccount(id int64) (*Account, error) {
	a := &Account{}
	//直接操作ID的简便方法
	has, err := x.Id(id).Get(a)
	//判断操作是否发生错误或对象是否存在
	if err != nil {
		return nil, err
	}else if !has {
		return nil, errors.New("Accoun does not exsit")
	}
	return a, nil
}
//用户存款
func makeDeposit(id int64, deposit float64) (*Account, error) {
	a, err := getAccount(id)
	if err != nil {
		return nil, err
	}
	a.Balance += deposit
	//对已有记录进行更新
	_, err1 := x.Update(a)
	return a, err1
}
//用户取款
func makeWithdraw(id int64, withdraw float64) (*Account, error) {
	a, err := getAccount(id)
	if err != nil {
		return nil, err
	}
	if a.Balance < withdraw {
		return nil, errors.New("Not enough balance")
	}
	a.Balance -= withdraw
	_, err1 := x.Update(a)
	return a, err1
}
//用户转账
func makeTransfer(id1 int64, id2 int64, balance float64) error {
	a1, err := getAccount(id1)
	if err != nil {
		return err
	}
	a2, err := getAccount(id2)
	if err != nil {
		return err
	}
	if a1.Balance < balance {
		return errors.New("Not enough balance")
	}
	
	//下面代码存在问题， 需要采用事务回滚来改进
	a1.Balance -= balance
	a2.Balance += balance
	//if _, err = x.Update(a1); err != nil {
	//	return err
	//}else if _, err = x.Update(a2); err != nil {
	//	return err
	//}
	//return nil
	sess := x.NewSession()
	defer sess.Close()
	//启动事务
	if err = sess.Begin(); err != nil {
		return err
	}
	if _, err = sess.Update(a1); err != nil {
		//发生错误时进行回滚
		sess.Rollback()
		return err
	}else if _, err = sess.Update(a2); err != nil {
		sess.Rollback()
		return err
	}
	return sess.Commit()
}
//按照ID正序排序返回所有账户
func getAccountAscId() (as []Account, err error) {
	//使用Find方法批量获取记录
	err = x.Find(&as)
	return as, err
}
//按照存款倒序排序返回所有账户
func getAccountDescBalance() (as []Account, err error) {
	//使用Desc方法使结果呈倒序排序
	err = x.Desc("balance").Find(&as)
	return as, err
}
//删除账户
func deleteAccount(id int64) error {
	//通过Delete方法删除记录
	_, err := x.Delete(&Account{Id:id})
	return err
}