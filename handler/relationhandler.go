package handler

import (
	"net/http"

	"github.com/Coreychen4444/Lite_TikTok/service"
	"github.com/gin-gonic/gin"
)

type RelationHandler struct {
	s *service.RelationService
}

func NewRelationHandler(s *service.RelationService) *RelationHandler {
	return &RelationHandler{s: s}
}

// 关注或取消关注
func (h *RelationHandler) FollowOrCancel(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	err := h.s.FollowOrCancel(token, to_user_id, action_type)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效,请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "关注失败" || err.Error() == "取消关注失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	if action_type == "1" {
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "关注成功"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "已取消关注"})
}

// 获取用户关注列表
func (h *RelationHandler) GetFollowings(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	followings, err := h.s.GetFollowings(token, user_id)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效,请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "获取关注列表失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "user_list": nil})
		return
	}
	if len(followings) == 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "暂无关注", "user_list": followings})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取关注列表成功", "user_list": followings})
}

// 获取用户粉丝列表
func (h *RelationHandler) GetFollowers(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	followers, err := h.s.GetFollowers(token, user_id)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效,请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "获取粉丝列表失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "user_list": nil})
		return
	}
	if len(followers) == 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "暂无粉丝", "user_list": followers})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取粉丝列表成功", "user_list": followers})
}

// 获取用户好友列表
func (h *RelationHandler) GetFriends(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	friends, err := h.s.GetFriends(token, user_id)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效,请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "获取好友列表失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "user_list": nil})
		return
	}
	if len(friends) == 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "暂无好友", "user_list": friends})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取好友列表成功", "user_list": friends})
}
