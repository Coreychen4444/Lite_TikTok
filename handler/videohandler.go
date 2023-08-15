package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/service"
	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	s *service.VideoService
}

func NewVideoHandler(s *service.VideoService) *VideoHandler {
	return &VideoHandler{s: s}
}

type GetVideoFlowRequest struct {
	LatestTime *string `json:"latest_time,omitempty"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	Token      *string `json:"token,omitempty"`       // 用户登录状态下设置
}

type GetVideoFlowResponse struct {
	NextTime   *int64        `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64         `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string       `json:"status_msg"`  // 返回状态描述
	VideoList  []model.Video `json:"video_list"`  // 视频列表
}

func (h *VideoHandler) GetVideoFlow(c *gin.Context) {
	var req GetVideoFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error(), "video_list": nil, "next_time": nil})
		return
	}
	if req.LatestTime == nil {
		latestTime := strconv.FormatInt(time.Now().Unix(), 10)
		req.LatestTime = &latestTime
	}
	if req.Token == nil {
		token := ""
		req.Token = &token
	}
	resp := &GetVideoFlowResponse{
		StatusCode: 1,
		StatusMsg:  nil,
		VideoList:  nil,
		NextTime:   nil,
	}
	video, err := h.s.GetVideoFlow(*req.LatestTime, *req.Token)
	if err != nil {
		respCode := http.StatusBadRequest
		errMsg := err.Error()
		if errMsg == "获取视频失败" {
			respCode = http.StatusInternalServerError
		}
		resp.StatusMsg = &errMsg
		c.JSON(respCode, resp)
		return
	}
	resp.StatusCode = 0
	if len(video) == 0 {
		resMsg := "没有更多视频了"
		resp.StatusMsg = &resMsg
		resp.VideoList = nil
		resp.NextTime = nil
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.VideoList = video
	nextTime := video[len(video)-1].PublishedAt
	resp.NextTime = &nextTime
	resp.StatusMsg = nil
	c.JSON(http.StatusOK, resp)
}

// 发布视频

func (h *VideoHandler) PublishVideo(c *gin.Context) {
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": "获取视频文件失败"})
		return
	}
	title := c.PostForm("title")
	token := c.PostForm("token")
	err = h.s.PublishVideo(file, title, token)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效,请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "创建视频失败" || err.Error() == "保存视频失败" || err.Error() == "生成封面失败" ||
			err.Error() == "发布视频失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "发布视频成功"})
}

type LikeVideoRequest struct {
	ActionType string `json:"action_type"` // 1-点赞，2-取消点赞
	Token      string `json:"token"`       // 用户鉴权token
	VideoID    string `json:"video_id"`    // 视频id
}

// 点赞视频和取消点赞
func (h *VideoHandler) LikeVideo(c *gin.Context) {
	var req LikeVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	err := h.s.LikeVideo(req.Token, req.VideoID, req.ActionType)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效, 请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "视频出错" || err.Error() == "请求参数错误" {
			respCode = http.StatusBadRequest
		} else if err.Error() == "点赞视频失败" || err.Error() == "取消点赞失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "操作成功"})
}

// 获取用户点赞的视频列表
func (h *VideoHandler) GetUserLike(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	videos, err := h.s.GetUserLike(token, user_id)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效, 请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "用户id格式错误" {
			respCode = http.StatusBadRequest
		} else if err.Error() == "获取列表失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "video_list": nil})
		return
	}
	if len(videos) == 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "该用户没有点赞的视频", "video_list": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取列表成功", "video_list": videos})
}

// 评论视频和删除评论
type CommentRequest struct {
	ActionType  string  `json:"action_type"`            // 1-发布评论，2-删除评论
	CommentID   *string `json:"comment_id,omitempty"`   // 要删除的评论id，在action_type=2的时候使用
	CommentText *string `json:"comment_text,omitempty"` // 用户填写的评论内容，在action_type=1的时候使用
	Token       string  `json:"token"`                  // 用户鉴权token
	VideoID     string  `json:"video_id"`               // 视频id
}

func (h *VideoHandler) CommentVideo(c *gin.Context) {
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	if req.ActionType == "1" {
		comment, err := h.s.CommentVideo(req.Token, req.VideoID, req.CommentText)
		if err != nil {
			respCode := http.StatusBadRequest
			if err.Error() == "token无效, 请重新登录" {
				respCode = http.StatusUnauthorized
			} else if err.Error() == "视频出错" || err.Error() == "请求参数错误" {
				respCode = http.StatusBadRequest
			} else if err.Error() == "评论视频失败" {
				respCode = http.StatusInternalServerError
			} else if err.Error() == "评论发表成功, 但返回评论信息失败" {
				c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": err.Error(), "comment": nil})
				return
			}
			c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "评论成功", "comment": comment})
	}
	if req.ActionType == "2" {
		err := h.s.DeleteComment(req.Token, *req.CommentID)
		if err != nil {
			respCode := http.StatusBadRequest
			if err.Error() == "token无效, 请重新登录" {
				respCode = http.StatusUnauthorized
			} else if err.Error() == "视频出错" || err.Error() == "请求参数错误" {
				respCode = http.StatusBadRequest
			} else if err.Error() == "删除评论失败" {
				respCode = http.StatusInternalServerError
			}
			c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "删除评论成功"})
	}
}

// 获取视频评论列表
func (h *VideoHandler) GetVideoComment(c *gin.Context) {
	video_id := c.Query("video_id")
	token := c.Query("token")
	comments, err := h.s.GetCommentList(token, video_id)
	if err != nil {
		respCode := http.StatusBadRequest
		if err.Error() == "token无效, 请重新登录" {
			respCode = http.StatusUnauthorized
		} else if err.Error() == "视频出错" {
			respCode = http.StatusBadRequest
		} else if err.Error() == "获取评论列表失败" {
			respCode = http.StatusInternalServerError
		}
		c.JSON(respCode, gin.H{"status_code": 1, "status_msg": err.Error(), "comment_list": nil})
		return
	}
	if len(comments) == 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "该视频没有评论", "comment_list": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取视频评论成功", "comment_list": comments})
}
