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

type FollowOrCancelReq struct {
	ActionType string `json:"action_type"` // 1-关注，2-取消关注
	ToUserID   string `json:"to_user_id"`  // 对方用户id
	Token      string `json:"token"`       // 用户鉴权token
}

// 关注或取消关注
func (h *RelationHandler) FollowOrCancel(c *gin.Context) {
	var req FollowOrCancelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	err := h.s.FollowOrCancel(req.Token, req.ToUserID, req.ActionType)
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
	if req.ActionType == "1" {
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
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "暂无关注", "user_list": nil})
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
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "暂无粉丝", "user_list": nil})
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
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "暂无好友", "user_list": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取好友列表成功", "user_list": friends})
}
