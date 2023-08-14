package repository

import (
	"github.com/Coreychen4444/Lite_TikTok/model"
	"gorm.io/gorm"
)

// 关注
func (r *DbRepository) CreateFollow(authorID, fansID int64) error {
	relation := model.Relation{
		AuthorID: authorID,
		FansID:   fansID,
	}
	tx := r.db.Begin()
	// 添加关注记录
	if err := r.db.Create(&relation).Error; err != nil {
		return err
	}
	// 更新作者的粉丝数
	if err := r.db.Model(&model.User{}).Where("id = ?", authorID).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
		return err
	}
	// 更新用户的关注数
	if err := r.db.Model(&model.User{}).Where("id = ?", fansID).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

// 取消关注
func (r *DbRepository) DeleteFollow(authorID, fansID int64) error {
	var relation model.Relation
	tx := r.db.Begin()
	// 删除关注记录
	if err := r.db.Where("author_id = ? and fans_id = ?", authorID, fansID).Delete(&relation).Error; err != nil {
		return err
	}
	// 更新作者的粉丝数
	if err := r.db.Model(&model.User{}).Where("id = ?", authorID).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
		return err
	}
	// 更新用户的关注数
	if err := r.db.Model(&model.User{}).Where("id = ?", fansID).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

// 获取用户关注列表
func (r *DbRepository) GetFollowList(userID int64) ([]model.User, error) {
	var following_id []int64
	err := r.db.Model(&model.Relation{}).Where("fans_id = ?", userID).Pluck("author_id", &following_id).Error
	if err != nil {
		return nil, err
	}
	var followings []model.User
	err = r.db.Where("id in (?)", following_id).Find(&followings).Error
	if err != nil {
		return nil, err
	}
	return followings, nil
}

// 获取用户粉丝列表
func (r *DbRepository) GetFansList(userID int64) ([]model.User, error) {
	var follower_id []int64
	err := r.db.Model(&model.Relation{}).Where("author_id = ?", userID).Pluck("fans_id", &follower_id).Error
	if err != nil {
		return nil, err
	}
	var followers []model.User
	err = r.db.Where("id in (?)", follower_id).Find(&followers).Error
	if err != nil {
		return nil, err
	}
	return followers, nil
}

// 获取用户好友列表
func (r *DbRepository) GetFriendList(userID int64) ([]model.User, error) {
	var following, follower []int64
	err := r.db.Model(&model.Relation{}).Where("fans_id = ?", userID).Pluck("author_id", &following).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&model.Relation{}).Where("author_id = ?", userID).Pluck("fans_id", &follower).Error
	if err != nil {
		return nil, err
	}
	friendID := GetFriID(following, follower)
	var friends []model.User
	err = r.db.Where("id in (?)", friendID).Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// 获取好友ID
func GetFriID(following, follower []int64) []int64 {
	var friendID []int64
	friendmap := make(map[int64]bool)
	for _, v := range following {
		friendmap[v] = true
	}
	for _, v := range follower {
		if friendmap[v] {
			friendID = append(friendID, v)
		}
	}
	return friendID
}
