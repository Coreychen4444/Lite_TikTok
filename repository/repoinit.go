package repository

import (
	"context"
	"log"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化mysql
func InitMysql() *gorm.DB {
	// 连接数据库(用户名和密码自己改)
	dsn := "root:44447777@tcp(127.0.0.1:3306)/tiktok_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error() + ", failed to connect database")
	}
	// 自动迁移
	err = db.AutoMigrate(&model.User{}, &model.Video{}, &model.VideoLike{}, &model.Comment{}, &model.Relation{}, &model.Message{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}
	log.Println("成功连接mysql数据库!")
	return db
}

// 初始化redis
func InitRedis() *redis.Client {
	// 连接redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "44447777",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis: " + err.Error())
	}
	log.Println("成功连接redis!")
	return rdb
}
