package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	s *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{s: s}
}

type RegisterResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

// 处理注册请求
func (h *UserHandler) Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	user_id, token, err := h.s.Register(username, password)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "生成token时出错" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "token": "", "user_id": -1})
		return
	}
	resp := &RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "注册成功",
		Token:      token,
		UserID:     user_id,
	}
	c.JSON(http.StatusOK, resp)
}

// 处理登录请求
func (h *UserHandler) Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	user_id, token, err := h.s.Login(username, password)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "生成token时出错" || err.Error() == "验证密码时出错" || err.Error() == "查找用户时出错" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "token": "", "user_id": -1})
		return
	}
	resp := &RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		Token:      token,
		UserID:     user_id,
	}
	c.JSON(http.StatusOK, resp)
}

type GetUserInfoResponse struct {
	StatusCode int64       `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string     `json:"status_msg"`  // 返回状态描述
	User       *model.User `json:"user"`        // 用户信息
}

// 处理获取用户信息请求
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	id, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		err := fmt.Errorf("用户id格式错误: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error(), "user": nil})
		return
	}
	user, err := h.s.GetUserInfo(id, token)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效,请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "该用户不存在" {
			respCode = http.StatusNotFound
		} else if err.Error() == "查找用户时出错" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "user": nil})
		return
	}
	statusMsg := "获取用户信息成功"
	resp := &GetUserInfoResponse{
		StatusCode: 0,
		StatusMsg:  &statusMsg,
		User:       user,
	}
	c.JSON(http.StatusOK, resp)
}

// 处理获取用户视频列表请求
func (h *UserHandler) GetUserVideoList(c *gin.Context) {
	token := c.Query("token")
	userID := c.Query("user_id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": "用户id格式错误", "video_list": nil})
		return
	}
	video, err := h.s.GetUserVideoList(id, token)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效,请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "获取视频失败" {
			respCode = http.StatusInternalServerError
		} else if err.Error() == "该用户不存在" {
			respCode = http.StatusNotFound
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "video_list": nil})
		return
	}
	if len(video) == 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "该用户没有发布任何视频", "video_list": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取用户视频列表成功", "video_list": video})
}
