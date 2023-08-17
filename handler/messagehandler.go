package handler

import (
	"net/http"

	"github.com/Coreychen4444/Lite_TikTok/service"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	s *service.MessageService
}

func NewMessageHandler(s *service.MessageService) *MessageHandler {
	return &MessageHandler{s: s}
}

// 获取聊天记录
func (h *MessageHandler) GetChatMessages(c *gin.Context) {
	to_user_id := c.Query("to_user_id")
	token := c.Query("token")
	pre_msg_time := c.Query("pre_msg_time")
	if pre_msg_time == "" {
		pre_msg_time = "0"
	}
	messages, err := h.s.GetChatMessages(token, to_user_id, pre_msg_time)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效,请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "获取聊天记录失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "message_list": nil})
		return
	}
	if len(messages) == 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "已是最新消息", "message_list": messages})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取聊天记录成功", "message_list": messages})
}

type SendMessageRequest struct {
	ActionType string `json:"action_type"` // 1-发送消息
	Content    string `json:"content"`     // 消息内容
	ToUserID   string `json:"to_user_id"`  // 对方用户id
	Token      string `json:"token"`       // 用户鉴权token
}

// 发送消息
func (h *MessageHandler) SendMessage(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	content := c.Query("content")
	if action_type == "1" {
		err := h.s.SendMessage(token, to_user_id, content)
		if err != nil {
			respCode := http.StatusBadRequest
			if err.Error() == "token无效,请重新登录" {
				respCode = http.StatusUnauthorized
			} else if err.Error() == "发送消息失败" {
				respCode = http.StatusInternalServerError
			}
			c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "发送消息成功"})
	}
}
