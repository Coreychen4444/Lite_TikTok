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
	/* 	r.POST("/douyin/user/register/", userHandler.Register)
	   	r.POST("/douyin/user/login/", userHandler.Login)
	   	r.GET("/douyin/user/", userHandler.GetUserInfo) */
	user := r.Group("/douyin/user")
	{
		user.POST("/register/", userHandler.Register)
		user.POST("/login/", userHandler.Login)
		user.GET("/", userHandler.GetUserInfo)
	}
	videoService := service.NewVideoService(repo)
	videoHandler := handler.NewVideoHandler(videoService)
	r.GET("/douyin/feed", videoHandler.GetVideoFlow)
	r.POST("/douyin/publish/action/", videoHandler.PublishVideo)
	r.Static("/public", "./public")
	r.GET("/douyin/publish/list/", userHandler.GetUserVideoList)
	r.POST("/douyin/favorite/action/", videoHandler.LikeVideo)
	r.GET("/douyin/favorite/list/", videoHandler.GetUserLike)
	r.POST("/douyin/comment/action/", videoHandler.CommentVideo)
	r.GET("/douyin/comment/list/", videoHandler.GetVideoComment)
	relationService := service.NewRelationService(repo)
	relationHandler := handler.NewRelationHandler(relationService)
	relation := r.Group("/douyin/relation")
	{
		relation.POST("/action/", relationHandler.FollowOrCancel)
		relation.GET("/follow/list/", relationHandler.GetFollowings)
		relation.GET("/follower/list/", relationHandler.GetFollowers)
		relation.GET("/friend/list/", relationHandler.GetFriends)
	}
	messageService := service.NewMessageService(repo)
	messageHandler := handler.NewMessageHandler(messageService)
	r.GET("/douyin/message/chat/", messageHandler.GetChatMessages)
	r.POST("/douyin/message/action/", messageHandler.SendMessage)
	return r
}
