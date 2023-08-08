package service

import (
	"fmt"
	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/repository"
)

type VideoService struct {
	r *repository.DbRepository
}

func NewVideoService(r *repository.DbRepository) *VideoService {
	return &VideoService{r: r}
}

// 获取视频列表
func (s *VideoService) GetVideoFlow(latest_time, token string) ([]model.Video, error) {
	video, err := s.r.GetVideoList(latest_time)
	if err != nil {
		return nil, fmt.Errorf("获取视频失败")
	}
	_, err = VerifyToken(token)
	if err != nil {
		return video, nil
	}
	return video, nil
}
