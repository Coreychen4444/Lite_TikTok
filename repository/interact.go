package repository

import (
	"fmt"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"gorm.io/gorm"
)

// 点赞
func (r *DbRepository) LikeVideo(user_id int64, video_id int64) error {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var videoLike model.VideoLike
	videoLike.UserID = user_id
	videoLike.VideoID = video_id
	if err := r.db.Create(&videoLike).Error; err != nil {

		return err
	}
	// 更新视频点赞数
	if err := r.db.Model(&model.Video{}).Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
		return err
	}
	// 更新用户点赞数
	if err := r.db.Model(&model.User{}).Where("id = ?", user_id).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
		return err
	}
	// 更新用户获赞数
	var authorId int64
	if err := r.db.Model(&model.Video{}).Where("id = ?", video_id).Pluck("user_id", &authorId).Error; err != nil {
		return err
	}
	if err := r.db.Model(&model.User{}).Where("id = ?", authorId).Update("total_favorited", gorm.Expr("total_favorited + ?", 1)).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

// 取消点赞
func (r *DbRepository) DislikeVideo(user_id int64, video_id int64) error {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var videoLike model.VideoLike
	if err := r.db.Where("user_id = ? and video_id = ?", user_id, video_id).Delete(&videoLike).Error; err != nil {
		return err
	}
	// 更新视频点赞数
	if err := r.db.Model(&model.Video{}).Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
		return err
	}
	// 更新用户点赞数
	if err := r.db.Model(&model.User{}).Where("id = ?", user_id).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
		return err
	}
	// 更新用户获赞数
	var authorId int64
	if err := r.db.Model(&model.Video{}).Where("id = ?", video_id).Pluck("user_id", &authorId).Error; err != nil {
		return err
	}
	if err := r.db.Model(&model.User{}).Where("id = ?", authorId).Update("total_favorited", gorm.Expr("total_favorited - ?", 1)).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

// 获取用户点赞的视频列表
func (r *DbRepository) GetUserLike(user_id int64) ([]model.Video, error) {
	videoId, err := r.GetUserLikeId(user_id)
	if err != nil {
		return nil, err
	}
	var videos []model.Video
	err = r.db.Preload("Author").Where("id in (?)", videoId).Order("published_at desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 获取用户喜欢的视频ID列表
func (r *DbRepository) GetUserLikeId(user_id int64) ([]int64, error) {
	var videoId []int64
	err := r.db.Model(&model.VideoLike{}).Where("user_id = ?", user_id).Pluck("video_id", &videoId).Error
	if err != nil {
		return nil, err
	}
	return videoId, nil
}

// 发布评论
func (r *DbRepository) CreateComment(comment *model.Comment) (*model.Comment, error) {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	if err := r.db.Create(comment).Error; err != nil {
		return nil, err
	}
	// 更新视频评论数
	if err := r.db.Model(&model.Video{}).Where("id = ?", comment.Video_id).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	id := comment.ID
	err := r.db.Preload("User").Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, fmt.Errorf("评论发表成功, 但返回评论信息失败")
	}
	return comment, nil
}

// 删除评论
func (r *DbRepository) DeleteComment(comment_id int64) error {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var comment model.Comment
	if err := r.db.Where("id = ?", comment_id).Delete(&comment).Error; err != nil {
		return err
	}
	// 更新视频评论数
	if err := r.db.Model(&model.Video{}).Where("id = ?", comment.Video_id).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

// 获取视频评论列表
func (r *DbRepository) GetCommentList(video_id int64) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("User").Where("video_id = ?", video_id).Order("create_date desc").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
