package datasource

import (
	"High/conf"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"sync"
)

var masterInstance *xorm.Engine
var once sync.Once
// 单例模式
func InstanceDbMaster() *xorm.Engine {
	// 使用Once保证创建实例的方法永远只能运行一次,就算在并发状态下也一定只执行一次
	once.Do(func() {
		masterInstance = NewDbMaster()
	})
	return masterInstance
}

func NewDbMaster() *xorm.Engine {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.DbMaster.User,
		conf.DbMaster.Pwd,
		conf.DbMaster.Host,
		conf.DbMaster.Port,
		conf.DbMaster.DataBase)
	instance,err := xorm.NewEngine(conf.DriverName, sourceName)
	if err != nil {
		log.Fatal("dbHelper NewEngine error",err)
		return nil
	}
	// TODO：调试时打开数据库展示
	instance.ShowSQL(true)
	masterInstance = instance
	return masterInstance
}
