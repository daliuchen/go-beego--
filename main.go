package main

import (
	_ "beegoTest/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var redisConnection redis.Conn

func main() {
	beego.Run()
}
func init() {
	beego.AddFuncMap("currentTime", getCurrentTime)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/test?charset=utf8")
}

func getCurrentTime() (out string) {
	time := time.Now()
	return time.Format("2006-01-02 15:04:05")
}
