package main

import (
	"github.com/Coreychen4444/Lite_TikTok/repository"
	"github.com/Coreychen4444/Lite_TikTok/routers"
)

func main() {
	//初始化数据库
	db := repository.InitMysql()
	//初始化redis
	rdb := repository.InitRedis()
	// 注册路由
	r := routers.InitRouter(db, rdb)
	r.Run("127.0.0.1:8080")
}
