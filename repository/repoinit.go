package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化mysql
func InitMysql() *gorm.DB {
	// 连接数据库(用户名和密码自己改)
	//dsn := "root:1234567@tcp(:3306)/tiktok_db?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("user:1234567@tcp(%s:3306)/tiktok_db?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_HOST"))
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
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		Password: "",
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
