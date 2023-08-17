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

type GetVideoFlowResponse struct {
	NextTime   *int64        `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64         `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string       `json:"status_msg"`  // 返回状态描述
	VideoList  []model.Video `json:"video_list"`  // 视频列表
}

func (h *VideoHandler) GetVideoFlow(c *gin.Context) {
	latest_time := c.Query("latest_time")
	token := c.Query("token")
	if latest_time == "0" {
		latest_time = strconv.FormatInt(time.Now().Unix(), 10)
	}
	resp := &GetVideoFlowResponse{
		StatusCode: 1,
		StatusMsg:  nil,
		VideoList:  nil,
		NextTime:   nil,
	}
	video, err := h.s.GetVideoFlow(latest_time, token)
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
		resp.VideoList = video
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
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
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

// 点赞视频和取消点赞
func (h *VideoHandler) LikeVideo(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	err := h.s.LikeVideo(token, video_id, action_type)
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
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "该用户没有点赞的视频", "video_list": videos})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取列表成功", "video_list": videos})
}

// 评论视频和删除评论
func (h *VideoHandler) CommentVideo(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	if action_type == "1" {
		comment_text := c.Query("comment_text")
		comment, err := h.s.CommentVideo(token, video_id, comment_text)
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
	if action_type == "2" {
		comment_id := c.Query("comment_id")
		err := h.s.DeleteComment(token, comment_id)
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
		c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "该视频没有评论", "comment_list": comments})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "获取视频评论成功", "comment_list": comments})
}
