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

type RegisterRequest struct {
	Password string `json:"password"` // 密码，最长32个字符
	Username string `json:"username"` // 注册用户名，最长32个字符
}

type RegisterResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

// 处理注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	user, token, err := h.s.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	resp := &RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Token:      token,
		UserID:     user.ID,
	}
	c.JSON(http.StatusOK, resp)
}

// 处理登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	user, token, err := h.s.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	resp := &RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Token:      token,
		UserID:     user.ID,
	}
	c.JSON(http.StatusOK, resp)
}

type GetUserInfoRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type GetUserInfoResponse struct {
	StatusCode int64       `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string     `json:"status_msg"`  // 返回状态描述
	User       *model.User `json:"user"`        // 用户信息
}

// 处理获取用户信息请求
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	var req GetUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	id, converr := strconv.Atoi(req.UserID)
	if converr != nil {
		err := fmt.Errorf("用户id格式错误: %w", converr)
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	user, err := h.s.GetUserInfo(int64(id), req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
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
