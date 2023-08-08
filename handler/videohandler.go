package handler

import (
	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}
	if req.LatestTime == nil {
		latestTime := ""
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
		errMsg := err.Error()
		resp.StatusMsg = &errMsg
		c.JSON(http.StatusBadRequest, resp)
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
	nextTime := video[len(video)-1].PublishedAt.Unix()
	resp.NextTime = &nextTime
	resp.StatusMsg = nil
	c.JSON(http.StatusOK, resp)
}
