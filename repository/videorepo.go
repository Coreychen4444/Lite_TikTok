package repository

import (
	"github.com/Coreychen4444/Lite_TikTok/model"
	"gorm.io/gorm"
)

// 获取视频列表

func (r *DbRepository) GetVideoList(latest_time int64) ([]model.Video, error) {
	var videos []model.Video
	err := r.db.Preload("Author").Where("published_at < ?", latest_time).Order("published_at desc").Limit(10).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 创建视频
func (r *DbRepository) CreateVideo(video *model.Video) error {
	tx := r.db.Begin()
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		txErr = err
		return err
	}
	if err := r.db.Create(video).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户投稿数
	if err := r.db.Model(&model.User{}).Where("id = ?", video.UserID).Update("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 获取某一用户投稿的视频列表
func (r *DbRepository) GetVideoListByUserId(user_id int64) ([]model.Video, error) {
	_, err := r.GetUserById(user_id)
	if err != nil {
		return nil, err
	}
	var videos []model.Video
	err = r.db.Preload("Author").Where("user_id = ?", user_id).Order("published_at desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
