package routers

import (
	"github.com/Coreychen4444/Lite_TikTok/handler"
	"github.com/Coreychen4444/Lite_TikTok/repository"
	"github.com/Coreychen4444/Lite_TikTok/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitRouter initialize routing information
func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	repo := repository.NewDbRepository(db)
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)
	user := r.Group("/douyin/user")
	{
		user.POST("/register", userHandler.Register)
		user.POST("/login", userHandler.Login)
		user.GET("/", userHandler.GetUserInfo)
	}
	videoService := service.NewVideoService(repo)
	videoHandler := handler.NewVideoHandler(videoService)
	r.GET("/douyin/feed", videoHandler.GetVideoFlow)
	r.POST("/douyin/publish/action", videoHandler.PublishVideo)
	r.Static("/public", "../public")
	r.GET("/douyin/publish/list", userHandler.GetUserVideoList)
	r.POST("/douyin/favorite/action", videoHandler.LikeVideo)
	r.GET("/douyin/favorite/list", videoHandler.GetUserLike)
	r.POST("/douyin/comment/action", videoHandler.CommentVideo)
	return r
}
